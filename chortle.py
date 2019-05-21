import requests
import json
import os
from pyunpack import Archive
import shutil

url = 'https://chorus.fightthe.pw'
dir_path = os.path.dirname(os.path.realpath(__file__))

print("\n"* 100)


def searchPrompt():
    search(input("What song would you like to search for?\n"))


def download(song_number, data):
    file_names = list(data['songs'][song_number]['directLinks'].keys())
    links = list(data['songs'][song_number]['directLinks'].values())

    folder_name = data['songs'][song_number]['name'] + ' - ' + data['songs'][song_number]['artist']

    try:
        os.mkdir(os.path.join(directory + folder_name))
    except FileExistsError:
        overwrite = input("Song already exists. Overwrite? (y/n)\n")
        if overwrite == 'y':
            shutil.rmtree(directory + folder_name)
            os.mkdir(os.path.join(directory + folder_name))
        else:
            searchPrompt()
    except PermissionError:
        print("I didn't have permissions to write the file.\n This is probably because you didn't set the songs folder location in the songs_directory folder.")
        print("The current songs folder location is set to:\n" + directory)
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
        with open(directory + folder_name + '/' + name, 'wb') as f:
            f.write(file_data.content)

        if name == 'archive':
            print("Unzipping archive...")
            try:
                Archive(directory + folder_name + '/archive').extractall(directory + folder_name)
                print("Cleaning up...")
                os.remove(directory + folder_name + '/archive')
            except:
                print("An unknown error occurred with the unzipping process, cleaning up...")
                shutil.rmtree(directory + folder_name)

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

        print(str(x+1) + ':   ', end='')
        print(songs[x]['name'], end='')
        print(' ' * (50-len(songs[x]['name'])), end='')
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


with open(dir_path + '/songs_directory.txt') as f:
    file = f.read()
    directory = file
    if(directory[len(directory)-2] != '/'):
        directory = directory.rstrip() + '/'
    else:
        directory = directory.rstrip()

print("Note: Make sure that you put your clone hero songs directory in the file. This program WILL NOT work without it.\n")

search(input("Type in what song you would like to search for:\n"))
