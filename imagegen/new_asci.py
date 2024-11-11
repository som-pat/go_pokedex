import sys
from PIL import Image
import requests
from io import BytesIO

def convert_to_ascii(image_url, new_width=64, final_width=64):
    # Fetch the image from the URL
    response = requests.get(image_url)
    img = Image.open(BytesIO(response.content))

    # Convert image to grayscale
    img = img.convert("L")  # "L" mode is for grayscale

    # Resize for terminal display
    width, height = img.size
    aspect_ratio = height / width
    new_height = int(aspect_ratio * new_width * 0.55)
    img = img.resize((new_width, new_height))

    # ASCII characters used for mapping (from dark to light)
    ascii_chars = " .:-=+*#%@"

    # Build the ASCII art string
    ascii_str = ""
    for y in range(new_height):
        line = ""
        for x in range(new_width):
            brightness = img.getpixel((x, y))  # Get brightness (0-255)
            char_index = brightness * (len(ascii_chars) - 1) // 255
            char = ascii_chars[char_index]
            line += char
        ascii_str += line + "\n"

    # Downsample by taking every nth row and column to match `final_width`
    downsample_ratio = max(1, new_width // final_width)
    downsampled_ascii = "\n".join(
        line[::downsample_ratio] for line in ascii_str.splitlines()[::downsample_ratio]
    )

    return downsampled_ascii

if __name__ == "__main__":
    image_url = sys.argv[1]
    ascii_art = convert_to_ascii(image_url)
    print(ascii_art)
