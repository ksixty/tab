from PIL import Image, ImageFont, ImageDraw

def read_font(img, ch, kern, c=(192, 192, 192, 255)):
    data = img.load()
    for y in range(img.size[1]):
        for x in range(img.size[0]):
            if data[x, y] == (255, 255, 255, 255):
                data[x, y] = c
            elif data[x, y][-1] == 0:
                data[x, y] = (255, 0, 255, 255)
            elif data[x, y] == (0, 0, 0, 255):
                data[x, y] = (0, 0, 0, 0)

    font = {None: kern}
    x = 0
    for c in ch:
        x_start = x
        while x < img.size[0] and img.getpixel((x, 0)) != (255, 0, 255, 255):
            x += 1
        x_end = x
        while x < img.size[0] and img.getpixel((x, 0)) == (255, 0, 255, 255):
            x += 1
        font[c] = img.crop((x_start, 0, x_end, img.size[1]))
        font[c].load()
    return font


def write_with_font(img, pos, text, font):
    x = pos[0]
    for c in text:
        img.paste(font[c], (x, pos[1]), font[c])
        x += font[c].size[0] + font[None]


japfon = read_font(Image.open("fonts/japfon.png"), "0123456789:", 3)
japfon_s = read_font(Image.open("fonts/japfon-sm.png"), "0123456789 -+:°X", 1)
japfon_s_dark = read_font(Image.open("fonts/japfon-sm.png"), "0123456789-+:°X", 1, c=(30, 30, 30, 255))
japfon_xs = read_font(Image.open("fonts/japfon-xs.png"), "0123456789-", 1)
