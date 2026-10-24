import cv2
import numpy as np
import HandTrackingModule as htm
import autopy

wCam, hCam = 640, 480
wScreen, hScreen = autopy.screen.size()
detector = htm.handDetector(maxHands=1)
cap = cv2.VideoCapture(0)
cap.set(3, wCam)
cap.set(4, hCam)
frameR = 100
smoothening = 7
pLocX, pLocY = 0, 0
cLocX, cLocY = 0, 0

while True:
    sucess, img = cap.read()
    img = detector.findHands(img)
    lmList, bbox = detector.findPosition(img)
    if len(lmList) != 0:
        x1, y1 = lmList[8][1:]
        x2, y2 = lmList[12][1:]
        fingers = detector.fingersUp()
        cv2.rectangle(img, (frameR, frameR), (wCam - frameR, hCam - frameR), (255, 0, 255), 2)
        if fingers[1] == 1 and fingers[2] == 0:
            print("Moving Mode!!")
            x3 = np.interp(x1, (frameR, wCam - frameR), (0, wScreen))
            y3 = np.interp(y1, (frameR, hCam - frameR), (0, hScreen))
            cLocX = pLocX + (x3 - pLocX) / smoothening
            cLocY = pLocY + (y3 - pLocY) / smoothening
            autopy.mouse.move(wScreen - cLocX, cLocY)
            cv2.circle(img, (x1, y1), 15, (255, 0, 255), cv2.FILLED)
            pLocX, pLocY = cLocX, cLocY
        if fingers[1] == 1 and fingers[2] == 1:
            length, img, lineInfo = detector.findDistance(8, 12, img)
            if length < 40:
                print("Clicking Mode!!")
                cv2.circle(img, (lineInfo[4], lineInfo[5]), 15, (0, 255, 0), cv2.FILLED)
                autopy.mouse.click()
    cv2.imshow("JAK AI Mouse Prog", img)
    cv2.waitKey(1)
