# kuronuri
A way to hide to text in japan,  is **KURONURI culture**. This tool implementation that way by Golang.

![image](https://github.com/yasutakatou/kuronuri/blob/pics/kuronuri.gif)

When you take over operation from your senior, you may look legacy shell script from old.

```
#!/bin/bash

USERID=admin
PASSWORD=password1234

wget --post-data="userid=$USERID&$PASSWORD" http://192.168.0.1:8080
```

### very, very, very worst!!

This makes you want to encrypt a little bit, doesn't it?<br>
But, If you encrypt all text, anyone can't see what's inside.<br>
If possible, **you want to encrypt a part of text**. This tool help you!

# usage

**This tool encrypt or decrypt a part of text, and create temporary file and run script.**<br>
**So, secret keyword can't know anyone , and we can know script's inside.**<br>

The earlier script example.

```
kuronuri.exe -noRun -dst=output.sh (USERID=;PASSWORD=:hogepassword) target.sh
```

This tool encrypt the sentence **after the key word**.

```
#!/bin/bash

USERID=Et3dTygDs6uOdxYiAA9cHWvrNmOvTRs1g-TQ9OTyBsrs
PASSWORD=VbRl8yLqWAbXossvQRlJRgJavonBaINn4hJvy5mJJtqkAa0wqH7RBg==

wget --post-data="userid=$USERID&$PASSWORD" http://192.168.0.220:8080
```

A case of decryption, **")"** is opposite direction.

```
kuronuri.exe )USERID=;PASSWORD=:hogepassword( output.sh
```

ver easy! (Bob Ross said)<br>

You want to hide your local address too if you can by **another key word**.

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

Way to decode by **two keywords**.

```
kuronuri.exe -noRun )USERID=;PASSWORD=:hogepassword( )//:fugapassword( output.sh
```

")" is opposite direction again. ver easy!

# important

This tool encrypt or decrypt **short length key word by AES**.<br>
So, Not very strong for brute-force attack.<br>
If you notice hacking, you must change key word encrypted.

# install

If you want to put it under the path, you can use the following.

```
go get github.com/yasutakatou/kuronuri
```

If you want to create a binary and copy it yourself, use the following.

```
git clone https://github.com/yasutakatou/kuronuri
cd kuronuri
go build kuronuri.go
```

or download binary from [release page](https://github.com/yasutakatou/kuronuri/releases).
save binary file, copy to entryed execute path directory.

# uninstall

delete that binary.
del or rm command. (it's simple!)

# options

"(" and ")" used encrypt or decrypt and this tool provide another options.<br>

|option name|default value|detail|
|:---|:---|:---|
|-dst|(random string)|output file name|
|-noRun|false|no Run and no Delete script|
|-dry|false|no Run and no Create script|
|-wrap|(Windows) busybox.exe <br> (Linux) bash|wrapper command|
|-opt|(Windows) bash <br> (Linux) [empty]|wrapper command arg option|

*note) 
When run on windows, windows don't provide unix command in the standard.<br>
So, this tool is using "busybox.exe".*<br>

[BUSYBOX](https://busybox.net/)<br>

*Download "busybox.exe" and put same directory this tool.*<br>

# License

ICU License
