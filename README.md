# Email2Whatsapp

The email2whatsapp is an innovative program written in Golang that allows users to find WhatsApp usernames from email addresses. This tool consults multiple databases to provide the most accurate results possible.

Inspired by Martin Vigo’s [email2phonenumber](https://github.com/martinvigo/email2phonenumber) tool, email2whatsapp takes functionality a step further, focusing on the popular messaging platform WhatsApp.

### Demo
- The email used has an account on PayPal and Magalu.
![Demo email2whatsapp](https://github.com/dsonbaker/email2whatsapp/blob/main/videos/demo_email2whatsapp1080p30fps.gif?raw=true)
##### Here are some key features of email2whatsapp

- Multiple Database Query: email2whatsapp is not limited to a single database. It consults multiple databases to ensure users get the most accurate and comprehensive results possible.

- Written in Golang: Golang is known for its efficiency and superior performance, making email2whatsapp a fast and reliable tool.

- Easy to Use: Despite its complexity under the hood, email2whatsapp is designed to be user-friendly, even for non-technical users.
#### Leaked digits in the Brazilian version

| Magalu         | Paypal         | PagBank        | Meli           | Rappi |
|---             |---             |---             |---             |---    |
| (01)9234*-**** | (0*)9***1-2345 | (0*)9****-1234 | (**)9****-1234 | (**)9****-1234 

- Currently in Brazil, all start with the number '9'.

## Technique to detect if a WhatsApp number exists.
- Instead of using paid APIs, the tutorial uses the feature of importing contacts from Google Contacts (contacts.google.com) using the same account linked to your smartphone.
- Wait for 5-10 minutes (depending on the number of contacts) for the WhatsApp contacts list to update.
- Use the command email2whatsapp -whatsapp and log in.
- The command will generate a file named numberphone.txt, which corresponds to the quantity of valid phone numbers found.

#### Disclaimer
Please note that responsible use of this tool is essential. It’s important to respect individuals’ privacy and rights when using such tools.