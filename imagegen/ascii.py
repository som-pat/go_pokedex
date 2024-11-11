import sys
from io import BytesIO
import cv2
import numpy as np
from math import floor,ceil
import collections
import requests

def pull_web_image(image_url):
    response = requests.get(image_url)
    if response.status_code == 200:
        img = np.frombuffer(BytesIO(response.content).getvalue(),dtype=np.uint8)
        img = cv2.imdecode(img, cv2.IMREAD_COLOR)
        if img is not None:
            return img
        else:
            raise ValueError("Image could not be decoded by OpenCV")
    else:
        raise ValueError(f"Could not fetch image. Status code: {response.status_code}")

def image_dimension(img):
    if img.shape[0]>512 and img.shape[1]>512:
        img = cv2.resize(img, (1024,1024), interpolation=cv2.INTER_AREA)
    else:
        img = cv2.resize(img, (128,128), interpolation=cv2.INTER_AREA)
    
    return img


def lab_contrast_enhance(img, factor):
    lab_img = cv2.cvtColor(img, cv2.COLOR_BGR2LAB)
    clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8, 8))
    lab_img[:, :, 0] = cv2.multiply(lab_img[:, :, 0], factor) #1.182
    lab_img[:, :, 0] = clahe.apply(lab_img[:, :, 0])  
    final_img = cv2.cvtColor(lab_img, cv2.COLOR_LAB2BGR)
    return final_img

def up_down_scaling(img, block_size):

    down_scaling = cv2.resize(img,(img.shape[1] // block_size, img.shape[0] // block_size), interpolation=cv2.INTER_AREA)    

    up_scaling = cv2.resize(down_scaling, (down_scaling.shape[1] * block_size, down_scaling.shape[0] * block_size), 
                            interpolation=cv2.INTER_NEAREST)

    return up_scaling


def sharpen(img):
    kernel = np.array([[-1,-1,-1],[-1,9,-1],[-1,-1,-1]])
    img = cv2.filter2D(img, -1, kernel)
    return img


def desat_graysc(img,cond):

    if cond==0:
        grayscale_image = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        return grayscale_image
    
    elif cond ==1:
        b = img[:,:,0].astype(np.float32)
        g = img[:,:,1].astype(np.float32)
        r = img[:,:,2].astype(np.float32)       

        desat = (0.299 * r + 0.587 * g + 0.114 * b).astype(np.uint8)
        return desat

    elif cond==2:
        b = img[:,:,0].astype(np.float32)
        g = img[:,:,1].astype(np.float32)
        r = img[:,:,2].astype(np.float32)
        desat = (0.299 * r + 0.587 * g + 0.114 * b).astype(np.uint8)
        grayscale_image = cv2.equalizeHist(desat)
        return grayscale_image
    
    elif cond==3:
        b = img[:,:,0].astype(np.float32)
        g = img[:,:,1].astype(np.float32)
        r = img[:,:,2].astype(np.float32)

        desat = (0.2126 * r + 0.7152 *g + 0.0722 *b).astype(np.uint8)
        return desat
    
    elif cond == 4:
        b = img[:,:,0].astype(np.float32)
        g = img[:,:,1].astype(np.float32)
        r = img[:,:,2].astype(np.float32)

        desat = np.sqrt( 0.299*r **2 + 0.587*g**2 + 0.114*b**2 ).astype(np.uint8)
        return desat


def enhance_edges(image,saturation, value, lightness):
    # increase saturation for bright edge 
    hsv_img = cv2.cvtColor(image, cv2.COLOR_BGR2HSV_FULL)
    hsv_img[:, :, 1] = cv2.multiply(hsv_img[:, :, 1], saturation)  
    hsv_img[:, :, 2] = cv2.multiply(hsv_img[:, :, 2], value)  
    enhanced_img = cv2.cvtColor(hsv_img, cv2.COLOR_HSV2BGR_FULL)
    
    # Apply CLAHE for enhance local contrast
    lab_img = cv2.cvtColor(enhanced_img, cv2.COLOR_BGR2LAB)
    clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8, 8))
    lab_img[:, :, 0] = cv2.multiply(lab_img[:, :, 0], lightness)
    lab_img[:, :, 0] = clahe.apply(lab_img[:, :, 0])  
    final_img = cv2.cvtColor(lab_img, cv2.COLOR_LAB2BGR)

    return final_img

def image_sharpen(img): #Image sharpening kernel
    kernel = np.array([[0,-1, 0], [-1,4,-1], [0,-1,0]]) 
    img = cv2.filter2D(img, -1, kernel)
    return img


def dog(img, kernel1, kernel2, sigma1, tau, th):    
    k = 1.6
    sigma2 = sigma1 * k 
    p = 11.5

    grad1 = cv2.GaussianBlur(img,(kernel1, kernel1),sigma1) 
    grad2 = cv2.GaussianBlur(img,(kernel2 ,kernel2),sigma2) 
    

    dog_sigma_k = ((1+tau) * grad1) - (tau * grad2)
    dog_sigma_k = grad1 + p * dog_sigma_k

    if th < np.percentile(dog_sigma_k, 85):
        th = np.percentile(dog_sigma_k, 85)

    dog_img = np.where(dog_sigma_k > th, dog_sigma_k, 0 ).astype(np.uint8)    
    return dog_img



def block_histogram(img_block, block_thresh):
    flat_block = img_block.flatten()
    non_zero_values = flat_block[flat_block > 0]

    if len(non_zero_values) < block_thresh:
        return 0

    median_angle = int(np.mean(non_zero_values))
    img_block = np.where(img_block > 0, median_angle, 0)    
    return img_block

def sobel_filter(img):
    Sx = cv2.Sobel(img, cv2.CV_32F, 1, 0, ksize=3)
    Sy = cv2.Sobel(img, cv2.CV_32F, 0, 1, ksize=3)

    grad_theta = np.arctan2(Sy, Sx)
    magnitude = np.sqrt(Sx **2 + Sy**2)

    return grad_theta, magnitude


def gradient_direction(edge, edge_Threshold, block_threshold = 10):
    grad_theta, magnitude = sobel_filter(edge)
    
    mag_edges = np.where(magnitude >= edge_Threshold, 255, 0).astype(np.uint8)
    
    grad_theta = np.degrees(grad_theta)
    grad_theta = np.round(grad_theta)
    new_grad = np.where(grad_theta < 0, ((180 - (grad_theta))%180) , grad_theta)
    new_grad = np.round(new_grad)
    grad_theta = np.round(grad_theta)
   

    char_size = (8,8)
    orient_map = np.zeros_like(edge)
    for i in range(0, mag_edges.shape[0], char_size[0]):
        for j in range(0, mag_edges.shape[1], char_size[1]):
            
            mag_block = mag_edges[i:i + char_size[0], j:j + char_size[1]]
            grad_block = grad_theta[i:i + char_size[0], j:j + char_size[1]]
            new_grad_block = new_grad[i:i + char_size[0], j:j + char_size[1]]
            
            
            if np.max(mag_block) !=0:              
                new_grad_block = np.where(mag_block == 0, 0, new_grad_block)
                array_max_block = block_histogram(new_grad_block,block_threshold)
                
                orient_map[i:i + char_size[0], j:j + char_size[1]] = array_max_block
                
            else:
                orient_map[i:i + char_size[0], j:j + char_size[1]] = new_grad_block
            
    return orient_map


def fetch_ascii_char(ascii_image, char_index, ascii_len, char_size=(8, 8)):
    
    x = (char_index % ascii_len) * char_size[0]  # Horizontal position
    y = (char_index // ascii_len) * char_size[1]  # Vertical position    
    # Crop ASCII character 
    ascii_char = ascii_image[y:y + char_size[1], x:x + char_size[0]]
    
    return ascii_char

def luminance_to_ascii_index(luminance, num_buckets=10):    
    return  floor((luminance / 255) * (num_buckets - 1))


def process_image_ascii(img):
    ascii_img = cv2.imread('ascii_py/ASCII_inside.png',cv2.IMREAD_GRAYSCALE) #grayscale for (8,80,3) to (8,80)
    char_size = (8,8)
    global im_count     
    ascii_len = 10
    dicte = {}

    ascii_art_image = np.zeros_like(img) #Empty board 
    

    for i in range(0, img.shape[0],char_size[0]):
        for j in range(0, img.shape[1],char_size[1]):   
            block = img[i:i + char_size[0], j:j + char_size[1]]
            
            if block.shape[0] != 8 or block.shape[1] != 8:
                continue
            
              
            average_luminance = np.mean(block)        

            ascii_index = luminance_to_ascii_index(average_luminance)
            dicte[ascii_index] = dicte.get(ascii_index,0)+1  

            ascii_char = fetch_ascii_char(ascii_img, ascii_index, ascii_len, char_size)      
            
            ascii_art_image[i:i + char_size[0], j:j + char_size[1]] = ascii_char
            
    

    print(dicte)
    return ascii_art_image




def edge_char_mapping(edge_ascii, edge_index, ascii_len ,char_size):
    x = (edge_index % ascii_len)  * char_size[0] # 0-40
    y = 0
    ascii_char = edge_ascii[y:y + char_size[1], x:x + char_size[0]]
    return ascii_char

def angle_to_ascii_index(angle):
    if angle == 0: return 0
    elif (angle>0 and angle <=23) or (angle>=157 and angle<=180): return 1
    elif (angle>=24 and angle <=78): return 4
    elif (angle>=79 and angle <=108): return 2
    elif (angle>=109 and angle<=156): return 3 


def ascii_edge_mapping(edge):
    edge_ascii=cv2.imread('ascii_py/edgesASCII.png',cv2.IMREAD_GRAYSCALE)
    char_size = (8,8)
    global ed_count
    ascii_len = 5

    edge_map = np.zeros_like(edge)
    dicte = {}
    for i in range(0,edge.shape[0],char_size[0]):
        for j in range(0, edge.shape[1],char_size[1]):

            edge_block = edge[i:i+char_size[1], j:j+char_size[0]]
            average_angle = np.max(edge_block)
            
            edge_index = angle_to_ascii_index(average_angle)
            if edge_index is None:
                print('edge_index',edge_index)
            dicte[edge_index] = dicte.get(edge_index,0)+1
            edge_char = edge_char_mapping(edge_ascii,edge_index, ascii_len ,char_size)
            edge_map[i:i+char_size[0],j:j+char_size[1]] = edge_char

    return edge_map


def overlay_images(img1, img2):    
    
    _, mask = cv2.threshold(img2, 50, 255, cv2.THRESH_BINARY)
    mask_inv = cv2.bitwise_not(mask)
    edges_from_img2 = cv2.bitwise_and(img2, img2, mask=mask)
    filled_from_img1 = cv2.bitwise_and(img1, img1, mask=mask_inv)
    combined_image = cv2.addWeighted(edges_from_img2, 0.85, filled_from_img1, 1, 0)
    
    return combined_image



red_count = 0

def mix_color_shader(ascii_img, original, gamma=1.2):

    ascii_norm = cv2.normalize(ascii_img.astype(np.float32), None, 0, 1, cv2.NORM_MINMAX)

    b, g, r = cv2.split(original)

    b = cv2.multiply(b.astype(np.float32), ascii_norm)
    g = cv2.multiply(g.astype(np.float32), ascii_norm)
    r = cv2.multiply(r.astype(np.float32), ascii_norm)

    colored_ascii = cv2.merge([b, g, r]).astype(np.uint8)

    inv_gamma = 1.0 / gamma
    lut = np.array([((i / 255.0) ** inv_gamma) * 255 for i in np.arange(0, 256)]).astype("uint8")
    cac= cv2.LUT(colored_ascii, lut)
    
    global red_count 
    red_count += 1
    file_loc = 'ascii_py/results/poke'+str(red_count)+'.jpg'
    cv2.imwrite(file_loc, cac) 


    return cac

# def image_to_ascii_text(img, width=16, ascii_chars="@%#*+=-:. "):
#     # Convert to grayscale if it's colored
#     if len(img.shape) == 3:
#         img = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

#     # Resize for terminal display
#     aspect_ratio = img.shape[0] / img.shape[1]
#     height = int(aspect_ratio * width * 0.55)  # Adjust for aspect ratio
#     img = cv2.resize(img, (width, height))

#     # Map each pixel to an ASCII character
#     ascii_str = ""
#     for pixel in img.flatten():
#         ascii_str += ascii_chars[pixel // (256 // len(ascii_chars))]

#     # Format ASCII art
#     ascii_lines = [ascii_str[index: index + width] for index in range(0, len(ascii_str), width)]
#     return "\n".join(ascii_lines)

def colorized_ascii_art(img, width=64, ascii_chars="@%#*+=-:. "):
    # Resize the image to the specified width and maintain the aspect ratio
    aspect_ratio = img.shape[0] / img.shape[1]
    height = int(aspect_ratio * width * 0.5)
    img = cv2.resize(img, (width, height), interpolation=cv2.INTER_AREA)

    ascii_art = ""
    for y in range(img.shape[0]):
        line = ""
        for x in range(img.shape[1]):
            # Get the RGB values
            b, g, r = img[y, x]
            
            # Map the pixel brightness to an ASCII character
            brightness = (0.299 * r + 0.587 * g + 0.114 * b)  # Luminosity formula for brightness
            char = ascii_chars[int(brightness / (256 / len(ascii_chars)))]
            
            # ANSI escape code for the color
            line += f"\033[38;2;{r};{g};{b}m{char}\033[0m"  # RGB foreground color
        ascii_art += line + "\n"
    
    return ascii_art

if __name__ == "__main__":
    image_url = sys.argv[1]
    og = pull_web_image(image_url)
    og = image_dimension(og)

    img = lab_contrast_enhance(og, factor=1.052)
    img = cv2.fastNlMeansDenoising(img, None, 20, 7, 21)
    img = lab_contrast_enhance(img, factor= 1.082) 
    img = desat_graysc(img,1)
    img = cv2.medianBlur(img,9)
    img = cv2.GaussianBlur(img,(7,7),0)
    img = up_down_scaling(img, block_size= 8)      
    img = process_image_ascii(img)


    edge = enhance_edges(og,saturation=1.32, value=0.85, lightness=1.29)        
    edge = desat_graysc(edge,4)
    edge = cv2.medianBlur(edge,9)
    edge = image_sharpen(edge)
    edge = cv2.medianBlur(edge,3)        
    edge = dog(edge, kernel1=3, kernel2 = 21, sigma1 = 0.98, tau=0.916, th =80)
    edge = cv2.medianBlur(edge,3)
    edge = gradient_direction(edge, edge_Threshold = 50, block_threshold = 10)
    edge = ascii_edge_mapping(edge)
    

    combi = overlay_images(img,edge)
    color_combi =mix_color_shader(combi, og)
    string_ascii = colorized_ascii_art(color_combi, width=16)
    print(string_ascii)