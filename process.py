import imageio
import os
from tqdm import tqdm


writer = imageio.get_writer('animated.gif', fps=1)

for im in tqdm(sorted(os.listdir('Data'), key=lambda x: int(x.split('.')[0]))):
    writer.append_data(imageio.imread(os.path.join('Data', im)))
writer.close()
