# Email2Whatsapp

The email2whatsapp is an innovative program written in Golang that allows users to find WhatsApp usernames from email addresses. This tool consults multiple databases to provide the most accurate results possible.

Inspired by Martin Vigo’s [email2phonenumber](https://github.com/martinvigo/email2phonenumber) tool, email2whatsapp takes functionality a step further, focusing on the popular messaging platform WhatsApp.

---

### Demo
+ Full Demo: https://www.youtube.com/watch?v=fgcDI-uWlTE
+ The email used has an account on PayPal and Magalu.<br><br>
[![Demo email2whatsapp](https://github.com/dsonbaker/email2whatsapp/blob/main/videos/demo_email2whatsapp1080p30fps.gif?raw=true)](https://vimeo.com/894134684)

##### Here are some key features of email2whatsapp

- Multiple Database Query: email2whatsapp is not limited to a single database. It consults multiple databases to ensure users get the most accurate and comprehensive results possible.

- Written in Golang: Golang is known for its efficiency and superior performance, making email2whatsapp a fast and reliable tool.

- Easy to Use: Despite its complexity under the hood, email2whatsapp is designed to be user-friendly, even for non-technical users.
---
#### Leaked digits in the Brazilian version


| **WebSites**          | **Leak Digits**  |
|---                    |---                |   
| **Magalu**            | (01)9234*-****    |
| **PayPal**            | (0*)9***1-2345    |
| **PagBank**           | (0*)9****-1234    |
| **Meli**              | (**)9****-1234    |
| **Rappi**             | (**)9****-1234    |




- Currently in Brazil, all start with the number '9'.
---
## Usage

- Scrape websites for phone number digits.
    ```
    email2whatsapp -email target@gmail.com
    ```
- Search for numbers with WhatsApp (the users.csv file should have been imported into Google Contacts).
    > Manually check after about 10 minutes if the WhatsApp contact list has synchronized, searching for "user " in the WhatsApp contact list.
    ```
    email2whatsapp -whatsapp
    ```
- Uses brute force to detect if it's the correct number (you'll need to solve some captchas for it to work properly).
    > Meli is the only one among these sources that returns the email initials.
    ```
    email2whatsapp -bruteforce paypal
    ```
---
## Technique to detect if a WhatsApp number exists.
- Instead of using paid APIs, the tutorial uses the feature of importing contacts from Google Contacts (contacts.google.com) using the same account linked to your smartphone.
- Wait for 8-15 minutes (depending on the number of contacts) for the WhatsApp contacts list to update.
- Use the command `email2whatsapp -whatsapp` and log in.
- The command will generate a folder named `./numberphone/all-numbers.txt`, which corresponds to the quantity of valid phone numbers found.
- If you know the photo of the person who owns the email, check the folder `./numberphone/profile/`, where public photos of each number are stored.
- If you didn't find the person's photo, try using the file `./numberphone/numbers-withoutProfile.txt` combined with the feature `email2whatsapp -bruteforce`.
---
## Help
- email2whatsapp -bruteforce
    - meli
        - Mercado Livre links the initials of the email. If they match the target email, there's a high possibility it's the WhatsApp number.
    - paypal
        - PayPal will only return whether the user exists or not.
    - twitter
        - Twitter, the initial requests can link to the email. However, after a certain number of requests, it will only return whether the user exists or not.
> Note that all these sites have captcha verification, hence requiring human assistance for captcha resolution. Therefore, the fewer possibilities, the better the result will be.
---
#### Disclaimer
> Please note that responsible use of this tool is essential. It’s important to respect individuals’ privacy and rights when using such tools.