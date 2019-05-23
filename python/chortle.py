import requests
import json
import os
from pyunpack import Archive
import shutil

url = 'https://chorus.fightthe.pw'
dir_path = os.path.dirname(os.path.realpath(__file__))

print("\n" * 100)


def searchPrompt():
    search(input("What song would you like to search for?\n"))


def download(song_number, data):
    file_names = list(data['songs'][song_number]['directLinks'].keys())
    links = list(data['songs'][song_number]['directLinks'].values())

    folder_name = data['songs'][song_number]['name'] + ' - ' + data['songs'][song_number]['artist']

    try:
        os.mkdir(os.path.join(directory, folder_name))
    except FileExistsError:
        overwrite = input("Song already exists. Overwrite? (y/n)\n")
        if overwrite == 'y':
            shutil.rmtree(os.path.join(directory,folder_name))
            os.mkdir(os.path.join(directory, folder_name))
        else:
            searchPrompt()
    except PermissionError:
        print(
            "I didn't have permissions to write the file.\n "
            "This is probably because you didn't set the correct songs folder location in the songs_directory folder."
            "The current songs folder location is set to:\n" + directory
        )
        quit()

    file_amount = len(data['songs'][song_number]['directLinks'])

    for x in range(file_amount):
        if file_names[x] == 'chart':
            name = 'notes.chart'
        elif file_names[x] == 'ini':
            name = 'song.ini'
        else:
            name = file_names[x]

        print("Downloading " + name)

        file_data = requests.get(links[x])

        print("Writing " + name)
        with open(os.join(directory, folder_name, name), 'wb') as f:
            f.write(file_data.content)

        if name == 'archive':
            print("Unzipping archive...")
            Archive(os.path.join(directory, folder_name, 'archive')).extractall(os.path.join(directory, folder_name))
            print("Cleaning up...")
            os.remove(os.path.join(directory, folder_name, 'archive'))

    print("Successfully downloaded!")
    answer = input("Search for another song? y/n:\n")
    if answer == 'y':
        searchPrompt()
    else:
        quit()


def search(search_term):
    page = url + '/api/search?query=' + str(search_term).replace(' ', '+')
    response = requests.get(page)
    data = response.json()
    songs = {}
    for x in range(len(data['songs'])):
        songs[x] = data['songs'][x]

        print(str(x + 1) + ':   ', end='')
        print(songs[x]['name'], end='')
        print(' ' * (50 - len(songs[x]['name'])), end='')
        print(songs[x]['artist'], end='')
        print(' ' * (35 - len(songs[x]['artist'])), end='')
        print(songs[x]['charter'])

    while True:
        chosen_song = input("Which song would you like to download? (s to make a different search, q to quit)\n")
        try:
            chosen_song = int(chosen_song)
        except ValueError:
            if chosen_song == 'q':
                quit()
            elif chosen_song == 's':
                searchPrompt()
                break
            else:
                print("Please type in a number!")

        if int(chosen_song) <= len(songs):
            break

    download(chosen_song - 1, data)


# read songs directory
with open(os.path.join(dir_path, 'songs_directory.txt'), 'r') as f:
    file = f.read()
    directory = file.rstrip()  # strip hidden \n

    if (len(directory) <= 3):
        # directory isn't set
        print("It looks like your songs location is not set, would you like to set it right now? (y/n)\n")

        if input() == 'y':
            with open(os.path.join(dir_path, 'songs_directory.txt'), 'w') as f:
                # rewrite songs directory
                directory = str(input("Type in your songs directory:\n"))
                f.write(directory)
                print('\n' * 5)

print("Current songs directory set to " + directory +
      "\n If this looks incorrect, you can edit songs_directory.txt, and reload the program")

searchPrompt()
