package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

/* SPEC
Generate a random salt ( base on website name ) to make it unique
Store Pepper in file ( created by user )

hash = sale + url + pepper

DEFINITIONS

SALT
In cryptography, a salt is random data that is used as an additional input
to a one-way function that "hashes" data, a password or passphrase.
Salts are closely related to the concept of nonce. The primary function of
salts is to defend against dictionary attacks or against its hashed
equivalent, a pre-computed rainbow table attack.[1]

Username	Salt value	String to be hashed	Hashed value = SHA256 (Password + Salt value)
user1	E1F53135E559C253	password123+E1F53135E559C253	72AE25495A7981C40622D49F9A52E4F1565C90F048F59027BD9C8C8900D5C3D8
user2	84B03D034B409D4E	password123+84B03D034B409D4E	B4B6603ABC670967E99C7E7F1389E40CD16E78AD38EB1468EC2AA1E62B8BED3A

PEPPER
In cryptography, a pepper is a secret added to an input such as a password
prior to being hashed with a cryptographic hash function.
A pepper performs a similar role to a salt, but while a salt is stored
alongside the hashed output, a pepper is not. A pepper usually meets
one of two criteria:

The password is not stored, and the 8-byte (64-bit) pepper 44534C70C6883DE2 is stored in a secure location separate to the hashed values.

Username	String to be Hashed	Hashed Value = SHA256(Password + Pepper)
user1	password123+44534C70C6883DE2	D63E21DF3A2A6853C2DC675EDDD4259F3B78490A4988B49FF3DB7B2891B3B48D
user2	password123+44534C70C6883DE2	D63E21DF3A2A6853C2DC675EDDD4259F3B78490A4988B49FF3DB7B2891B3B48D
*/

func calc(input string) string {
	// check whether we want sha1 or sha256
	// The deciding factor is the existance of the file sha1 or sha256
	hash := CheckCryptoType()

	// get pepper from public.key file
	// only the first line is read
	// This can contain any text but make it unique
	pepper, err := GetPepper("public.key")
	if err != nil {
		panic("public.key missing")
	}

	// calculate the salt from the url
	// because we don't need to decrypt and we need a unique salt
	// for each url, the salt is based on the url
	salt, _ := GetSalt(hash, input)

	// add everything to form the text before hashing
	input = salt + input + pepper

	// generate a hash
	result := GenerateHash(hash, input)
	fmt.Printf("\tPlaintext = [%s]\n\tEncrypted = [%s] \n", input, result)
	return result
}

func CheckCryptoType() hash.Hash {
	if _, err := os.Stat("sha1"); err == nil {
		fmt.Println("Using sha1")
		return sha1.New()
	}
	if _, err := os.Stat("sha256"); err == nil {
		fmt.Println("Using sha256")
		return sha256.New()
	}

	fmt.Println("Using default crypto")
	return sha1.New()
}

func GetPepper(fileName string) (string, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return "", err
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		fmt.Printf("Read from %s\n", fileName)
		break
	}

	if err != io.EOF {
		panic(" > 1 line only required:\n")
	}

	return line, nil
}

func GetSalt(h hash.Hash, input string) (string, error) {
	h.Write([]byte(input))
	bs := h.Sum(nil)
	result := fmt.Sprintf("%x", bs)
	return result, nil
}

func GenerateHash(h hash.Hash, str string) string {
	h.Write([]byte(str))
	bs := h.Sum(nil)
	result := fmt.Sprintf("%x", bs)
	return result
}
