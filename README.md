# kuronuri
A way to hide to text in japan,  is KURONURI culture. This tool implementation that way by Golang.

![image](https://github.com/yasutakatou/kuronuri/blob/pics/kuronuri.gif)

When you take over operation from your senior, you may look legacy shell script from old.

```
#!/bin/bash

USERID=admin
PASSWORD=password1234

wget --post-data="userid=$USERID&$PASSWORD" http://192.168.0.1:8080
```

## very, very, very worst!!

This makes you want to encrypt a little bit, doesn't it?
But, If you encrypt all text, anyone don't know about text detail.
If possible, you want to encrypt a part of text. This tool help you!

# usage

The earlier script example.

```
kuronuri.exe -noRun -dst=output.sh (USERID=;PASSWORD=:hogepassword) target.sh
```

This tool encrypt the sentence after the keyword.

```
#!/bin/bash

USERID=Et3dTygDs6uOdxYiAA9cHWvrNmOvTRs1g-TQ9OTyBsrs
PASSWORD=VbRl8yLqWAbXossvQRlJRgJavonBaINn4hJvy5mJJtqkAa0wqH7RBg==

wget --post-data="userid=$USERID&$PASSWORD" http://192.168.0.220:8080
```

A case of decryption, ")" is opposite direction.

```
kuronuri.exe )USERID=;PASSWORD=:hogepassword( output.sh
```

ver easy! (Bob Ross say)

You want to hide your local address too if you can by another key word.

```
kuronuri.exe -noRun -dst=output.sh (USERID=;PASSWORD=:hogepassword) (//:fugapassword) auth.sh
```

You could.

```
#!/bin/bash

USERID=qDCHaNFLiHI7UtFM_r2_dZzGUacQnFLfXZ8gOw1gsoay
PASSWORD=T_EHdS2C6RAUhstj3PIi124uRwhkWgEE_A67Z7e4NUx9rtso2IX21g==

wget --post-data="userid=$USERID&$PASSWORD" http://eNU4qWl3qJ7lVxL14lDTrMXkyUQuuMXwzM9nsbowd8ZVDJFyC1XJ5KRXJsUvAA==
```

Your keyword way to decode.

```
kuronuri.exe -noRun )USERID=;PASSWORD=:hogepassword( )//:fugapassword( output.sh
```
