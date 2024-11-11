import sys
from PIL import Image
import requests
from io import BytesIO

def convert_to_ascii(image_url, new_width=64,final_width=16):
    # Fetch the image from the URL
    response = requests.get(image_url)
    img = Image.open(BytesIO(response.content))

    # Convert image to RGB to ensure we have color data
    img = img.convert("RGB")

    # Resize for terminal display
    width, height = img.size
    aspect_ratio = height / width
    new_height = int(aspect_ratio * new_width * 0.55)
    img = img.resize((new_width, new_height))

    # ASCII characters used for mapping (from dark to light)
    ascii_chars = " .:-=+*#%@"

    # Initialize ASCII art string
    ascii_str = ""
    pixel_count = 0
    for pixel in img.getdata():
        # Handle different pixel formats
        if isinstance(pixel, int):
            # Grayscale image
            r = g = b = pixel
        elif len(pixel) == 4:
            # RGBA image
            r, g, b, a = pixel
        elif len(pixel) == 3:
            # RGB image
            r, g, b = pixel
        else:
            raise ValueError(f"Unexpected pixel format: {pixel}")

        # Calculate brightness and map to ASCII character
        brightness = int(0.299 * r + 0.587 * g + 0.114 * b)  # Luminosity formula for brightness
        char_index = brightness * (len(ascii_chars) - 1) // 255
        char = ascii_chars[char_index]

        # Add ANSI escape code for color
        ascii_str += f"\033[38;2;{r};{g};{b}m{char}\033[0m"

        pixel_count += 1
        # Add newline after each row
        if pixel_count % new_width == 0:
            ascii_str += "\n"

    ascii_lines = ascii_str.splitlines()

    # Downsample the ASCII art to match the final desired width
    downsample_ratio = max(1, new_width // final_width)
    downsampled_ascii = "\n".join(
        line[::downsample_ratio] for line in ascii_lines[::downsample_ratio]
    )

    return downsampled_ascii


if __name__ == "__main__":
    image_url = sys.argv[1]
    ascii_art = convert_to_ascii(image_url)
    print(ascii_art)
