import math
import requests as internet
import time
from bs4 import BeautifulSoup as htmlparser

cookie = ""
userid = 0 #This is the user id of your account

excludeid = 514593823 #This is a value to exclude removing someone from your friends list. Limited to one until later update(I just made this for myself at first)



headers = {
    "authority": "auth.roblox.com",
    "accept": "application/json, text/plain, */*",
    "content-type": "application/json",
    "_gcl_au": "1.1.191582745.1703431727",
    "GuestData": f"UserID=-{1521315+(math.floor(1521315*1000))})"
}

session = internet.Session()
session.cookies[".ROBLOSECURITY"] = cookie

def getCSRF():
    session.cookies["x-csrf-token"] = None
    session.headers["x-csrf-token"] = None

    response = session.get("https://www.roblox.com/catalog/17218442369/")

    html = htmlparser(response.text, "html.parser")
    csrf_tag = html.find("meta", {"name": "csrf-token"})
    csrf_token = csrf_tag["data-token"]

    session.cookies["x-csrf-token"] = csrf_token
    session.headers["x-csrf-token"] = csrf_token
    session.headers["withCredentials"] = "true"

print("Devious Lick RF has started up.")
time.sleep(2.5)
print("Devious Lick RF has successfully started up! Removing friends...")
time.sleep(1)

def main():
    idList = []

    newreq = session.get(f"https://friends.roblox.com/v1/users/{userid}/friends")

    for req in newreq.json()["data"]:
        if req["id"] not in idList:
            idList.append(int(req["id"]))

    getCSRF()

    for id in idList:
        if id != excludeid:
            status_code = session.post(f"https://friends.roblox.com/v1/users/{id}/unfriend")
            print(status_code.status_code)
main()