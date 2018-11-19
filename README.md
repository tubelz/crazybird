# Crazy Bird

Game made to demonstrate how to create games using Macaw

Quick video: https://youtu.be/YRYtKJVqQCI

## Running with Docker

Make sure you have added in the .env file the three required variables: 
  - APP_NAME: this is your game's name;
  - REPOSITORY: path of your game;
  - DISPLAY: you can define a window or just use $DISPLAY to get the active window.

### Dev

If you want to run in dev mode, you have two options: using docker or docker-compose. Regardless of the option
you have picked, you must build the project within the container. After the project has been compiled, you can
run the application.

This option will mount a volume that will link to your project path

1. `docker-compose build`
2. `docker-compose run macaw`

### Production

Use this when there are no more changes to be applied to the version of the project. That way you'll have a static
image that you can utilize without worrying much.

The difference between this option and dev is that here we make a copy of the files and we don't use the bind volume

1. `docker-compose -f docker-compose.yml -f docker-compose.prod.yml build`
2. `docker-compose -f docker-compose.yml -f docker-compose.prod.yml run macaw`

## License
The code here is under the zlib license. You can read more [here](https://github.com/tubelz/crazybird/LICENSE.txt)
