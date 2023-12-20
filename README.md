# Email2Whatsapp

The email2whatsapp is an innovative program written in Golang that allows users to find WhatsApp usernames from email addresses. This tool consults multiple databases to provide the most accurate results possible.

Inspired by Martin Vigo’s [email2phonenumber](https://github.com/martinvigo/email2phonenumber) tool, email2whatsapp takes functionality a step further, focusing on the popular messaging platform WhatsApp.

---

### Demo
+ Full Demo: https://youtu.be/YcWO0BmVKK0
+ The email used has an account on PayPal and Magalu.<br><br>
![Demo email2whatsapp](https://github.com/dsonbaker/email2whatsapp/blob/main/videos/demo_email2whatsapp1080p30fps.gif?raw=true)

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
## Instalation
```
go install -v github.com/dsonbaker/email2whatsapp@latest
```
---
## Usage

- Scrape websites for phone number digits.
    ```
    email2whatsapp -email target@gmail.com
    ```
- Search for numbers with WhatsApp.
    > Connect your WhatsApp using the QR code.
    ```
    echo 5521912345678 | email2whatsapp -whatsapp
    ```
- Uses brute force to detect if it's the correct number (you'll need to solve some captchas for it to work properly).
    > Meli returns the email initials.
    ```
    echo 5521912345678 | email2whatsapp -bruteforce meli
    ```

- Uses brute force to detect if it's the correct number (No need to solve a captcha).
    > Microsoft returns the email initials.
    ```
    echo 5521912345678 | email2whatsapp -bruteforce microsoft
    ```
---
## Technique to detect if a WhatsApp number exists.
- The [Whatsmeow](https://github.com/tulir/whatsmeow) project was used to establish a connection with the WhatsApp protocol.
- The process of searching for valid numbers, in the case of a list of 5000 possible numbers, can take around 15 minutes.
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
    - google
        - Google will only return if the number is linked to an account.
    - microsoft
        - Microsoft will return some characters of the email linked to the number.
> Note that some of these websites have captcha verification, thus requiring human assistance for captcha resolution. Therefore, the fewer the possibilities, the better the outcome.
---
#### Disclaimer
> Please note that responsible use of this tool is essential. It’s important to respect individuals’ privacy and rights when using such tools.