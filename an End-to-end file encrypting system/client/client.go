package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"bytes"
	"encoding/json"
	"strconv"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation

	// Useful for formatting strings (e.g. `fmt.Sprintf`).

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

type UUID = uuid.UUID
type PKEEncKey = userlib.PKEEncKey
type PKEDecKey = userlib.PKEDecKey
type DSSignKey = userlib.DSSignKey
type DSVerifyKey = userlib.DSVerifyKey

type User struct {
	EncInnerUserStructPointer UUID
	HMACEncInnerUserStruct    []byte
	username                  string
	rootKey                   []byte
}

type InnerUser struct {
	//Username  string -- Do not want to worry about AppendToFile() scaling with username
	PKEDecKey PKEDecKey
	DSSignKey DSSignKey
}

/* Cycle: Sentinel <-> Outer File 1 <-> Outer File 2 <-> Sentinel <-> Outer File 1 <-> ..... etc */
type OuterFile struct {
	EncInnerFileStructPointer UUID
	HMACEncInnerFileStruct    []byte
}

type InnerFile struct {
	LinkedListNumber int
	Content          []byte
}

type OuterInvitation struct {
	EncInnerInvitationStructPointer UUID // Encrypted with public key of the recipient
	DigitalSignature                []byte
}

// type InnerInvitation struct {
// 	// Note: Changes made here so that there is no bandwidth scaling with anything at all
// 	PointerToOuterSentinelStructPointer UUID // Pointer to pointer to the Outer Sentinel File Struct for the file
// 	PointerToSymKey                     UUID // Pointer to pointer to symmetric key for the file
// 	PointerToHMACKey                    UUID // Pointer to pointer to HMAC for the file
// }

type InnerInvitation struct {
	PointerToEncSharedFileRootKey UUID
	//FileRootKey []byte
}

type RecipientUsernames struct {
	ListOfHashedRecipientUsernames [][]byte
}

////////////////////////////////////
// HELPER FUNCTIONS
////////////////////////////////////

func IsEmptyUsername(username string) (isEmpty bool) {
	return len(username) == 0
}

func CheckUserExists(username string, password string) (exists bool) {
	rootKey := userlib.Argon2Key([]byte(password), []byte(username), 16) //PBKDF
	hashedRootKeyString := string(userlib.Hash(rootKey))
	keystorePKEENCKeyLocation := hashedRootKeyString + "PKEEncKey"
	_, ok := userlib.KeystoreGet(keystorePKEENCKeyLocation)
	return ok
}

func CheckUserExistsNOPassword(rootKey []byte) (exist bool) {
	hashedRootKeyString := string(userlib.Hash(rootKey))
	keystorePKEENCKeyLocation := hashedRootKeyString + "PKEEncKey"
	_, ok := userlib.KeystoreGet(keystorePKEENCKeyLocation)
	return ok
}

func GetUserPublicKeys(rootKey []byte) (PKEEncKey PKEEncKey, DSVerifyKey DSVerifyKey, err error) {

	// Get the hashedRootKeyString
	hashedRootKeyString := string(userlib.Hash(rootKey))

	// Check if the Datastore has the public keys, and if it doesn't return error
	PKEEncKey, ok := userlib.KeystoreGet(hashedRootKeyString + "PKEEncKey")
	if !ok {
		return userlib.PKEEncKey{}, userlib.DSVerifyKey{}, errors.New("keystore does not contain the PKEEncKey")
	}

	DSVerifyKey, ok = userlib.KeystoreGet(hashedRootKeyString + "DSVerifyKey")
	if !ok {
		return userlib.PKEEncKey{}, userlib.DSVerifyKey{}, errors.New("keystore does not contain the DSVerifyKey")
	}

	// Return the values
	return PKEEncKey, DSVerifyKey, nil
}

func GetUserPublicKeysFromUsername(username string) (PKEEncKey PKEEncKey, DSVerifyKey DSVerifyKey, err error) {

	// Check if the Datastore has the public keys, and if it doesn't return error
	PKEEncKey, ok := userlib.KeystoreGet(username + "PKEEncKey")
	if !ok {
		return userlib.PKEEncKey{}, userlib.DSVerifyKey{}, errors.New("keystore does not contain the PKEEncKey")
	}

	DSVerifyKey, ok = userlib.KeystoreGet(username + "DSVerifyKey")
	if !ok {
		return userlib.PKEEncKey{}, userlib.DSVerifyKey{}, errors.New("keystore does not contain the DSVerifyKey")
	}

	// Return the values
	return PKEEncKey, DSVerifyKey, nil
}

// Get user's privates for signature
func GetUserPrivateKeys(rootKey []byte) (PKEDecKey PKEDecKey, DSSignKey DSSignKey, err error) {
	// Access Inner User Struct
	// Get user and then decrypt

	innerUserStructPtr, err := LoadInnerUserStructFromDatastore(rootKey)
	if err != nil {
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, err
	}

	PKEDecKey = (*innerUserStructPtr).PKEDecKey
	DSSignKey = (*innerUserStructPtr).DSSignKey

	return PKEDecKey, DSSignKey, err
}

/* Returns pointers to the inner and outer user structs. Also generates the Digital Signature and RSA Keys, stores them in Keystore */
func CreateAndStoreInnerAndOuterUserStructs(username string, password string) (innerUserStructPtr *InnerUser, outerUserStructPtr *User, err error) {

	rootKey := userlib.Argon2Key([]byte(password), []byte(username), 16) //PBKDF

	hashedRootKeyString := string(userlib.Hash(rootKey))

	//Digital Signature Keys
	DSSignKey, DSVerifyKey, err := userlib.DSKeyGen()
	if err != nil {
		return nil, nil, errors.New("an error occured in Signature Key Generation") //no error should ever happen here?
	}
	// Store the public keys twice, with different ways to access them (for convenience later)
	err = userlib.KeystoreSet(hashedRootKeyString+"DSVerifyKey", DSVerifyKey)
	if err != nil {
		return nil, nil, errors.New("could not set the DSVerifyKey in the Keystore")
	}
	err = userlib.KeystoreSet(username+"DSVerifyKey", DSVerifyKey)
	if err != nil {
		return nil, nil, errors.New("could not set the DSVerifyKey in the Keystore")
	}

	//RSA Keys
	PKEEncKey, PKEDecKey, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, nil, errors.New("an error occured in Public Key Generation") //no error should ever happen here?
	}
	err = userlib.KeystoreSet(hashedRootKeyString+"PKEEncKey", PKEEncKey)
	if err != nil {
		return nil, nil, errors.New("could not set the PKEEncKey in the Keystore")
	}
	err = userlib.KeystoreSet(username+"PKEEncKey", PKEEncKey)
	if err != nil {
		return nil, nil, errors.New("could not set the PKEEncKey in the Keystore")
	}

	// Step 4: Create the InnerUser struct
	innerUser := InnerUser{
		PKEDecKey: PKEDecKey,
		DSSignKey: DSSignKey,
	}

	// Step 5: Generate the root key, Outer User Struct Key, (Symmetric) Encryption Key, and HMAC Key for Inner User Struct
	outerUserStructKey, innerUserStructEncryptKey, innerUserStructHMACKey, err := GenerateUserStructKeys(rootKey)
	if err != nil {
		return nil, nil, err
	}

	// Step 6: Store the ENCRYPTED InnerUser struct in the Datastore
	datastorePointerToInnerUserStruct, err := GenerateUUID(innerUserStructEncryptKey)
	if err != nil {
		return nil, nil, err
	}

	_, encryptedInnerUserStruct, err := MarshalEncryptStoreInnerStruct(innerUser, innerUserStructEncryptKey, datastorePointerToInnerUserStruct)
	if err != nil {
		return nil, nil, err
	}

	// Step 7: Generate the (Outer) User Struct
	HMACEncryptedInnerUserStruct, err := userlib.HMACEval(innerUserStructHMACKey, encryptedInnerUserStruct)
	if err != nil {
		return nil, nil, err
	}

	outerUser := User{
		EncInnerUserStructPointer: datastorePointerToInnerUserStruct,
		HMACEncInnerUserStruct:    HMACEncryptedInnerUserStruct,
		username:                  username,
		rootKey:                   rootKey,
	}

	// Marshal and Store the OuterUserStruct in the DataStore
	datastorePointerToUserStruct, err := GenerateUUID(outerUserStructKey)
	if err != nil {
		return nil, nil, err
	}

	_, err = MarshalStoreOuterStruct(outerUser, datastorePointerToUserStruct)
	if err != nil {
		return nil, nil, err
	}

	return &innerUser, &outerUser, nil
}

/* Returns pointers to the inner and outer file structs. If linkedListNumber = 0, then we need to make the sentinel (file) structs, and nothing else. If linkedListNumber > 0, then we need to make the corresponding (file) structs, and nothing else */
func CreateAndStoreInnerAndOuterFileStructs(filename string, content []byte, fileRootKey []byte, revokeNumber int, linkedListNumber int, linkedListSize int) (err error) {

	/* GENERATING KEYS */

	_, innerFileStructEncryptKey, innerFileStructHMACKey, err := GenerateFileStructKeys(fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return err
	}

	// STEP 2: Generate Datastore UUID's (locations) for Outer + Inner File struct
	datastorePointerToOuterFileStruct, datastorePointerToInnerFileStruct, err := GetDatastorePointersToFileStructs(fileRootKey, filename, revokeNumber, linkedListNumber)
	if err != nil {
		return err
	}

	var innerFileStruct InnerFile
	var outerFileStruct OuterFile

	// There are two cases to check.
	// Case 1: linkedListNumber = 0 ----> create the sentinel structs, and assign sentinelInnerFileStruct.LinkedListNumber = linkedListSize
	// Case 2: linkedListNumber > 0 ----> create the corresponding file structs, and assign innerFileStruct.LinkedListNumber = linkedListNumber (the number doensn't really matter)

	if linkedListNumber == 0 {

		/* SENTINEL STRUCT SECTION */

		// STEP 5: Create the Sentinel InnerFile struct
		sentinelInnerFileStruct := InnerFile{
			LinkedListNumber: linkedListSize,
			Content:          nil,
		}

		innerFileStruct = sentinelInnerFileStruct

	} else if linkedListNumber > 0 {

		/* FILE STRUCT SECTION */

		// STEP 10: Create the InnerFile struct
		innerFileStruct = InnerFile{
			LinkedListNumber: linkedListNumber,
			Content:          content,
		}

	} else {
		return errors.New("invalid LinkedListNumber passed in")
	}

	// STEP 6: Marshal, Encrypt, then Store the Sentinel InnerFile struct to the Datastore
	_, encryptedInnerFileStruct, err := MarshalEncryptStoreInnerStruct(innerFileStruct, innerFileStructEncryptKey, datastorePointerToInnerFileStruct)
	if err != nil {
		return err
	}

	// STEP 7: Compute the HMAC of the (Encrypted) Sentinel InnerFile struct
	innerFileStructHMAC, err := userlib.HMACEval(innerFileStructHMACKey, encryptedInnerFileStruct)
	if err != nil {
		return errors.New("could not compute the HMAC of the Inner File Struct")
	}

	// STEP 8: Create the Sentinel OuterFile struct
	outerFileStruct = OuterFile{
		EncInnerFileStructPointer: datastorePointerToInnerFileStruct,
		HMACEncInnerFileStruct:    innerFileStructHMAC,
	}

	// STEP 9: Marshal then Store the Sentinel OuterFile struct in the Datastore
	_, err = MarshalStoreOuterStruct(outerFileStruct, datastorePointerToOuterFileStruct)
	if err != nil {
		return err
	}

	return nil
}

/* Returns err if failed at any point (including the HMAC check). Returns the pointer to the outer user struct */
func LoadOuterUserStructFromDatastore(rootKey []byte) (outerUserStructPtr *User, err error) {

	// Generate Keys
	outerUserStructKey, _, innerUserStructHMACKey, err := GenerateUserStructKeys(rootKey)
	if err != nil {
		return nil, err
	}

	// Get Datastore Location of (Outer) User Struct
	datastorePointerToUserStruct, err := GenerateUUID(outerUserStructKey)
	if err != nil {
		return nil, errors.New("could not generate UUID for Outer User Struct")
	}

	marshalledUserStruct, ok := userlib.DatastoreGet(datastorePointerToUserStruct)
	if !ok {
		return nil, errors.New("datastore does not Contain User Struct")
	}

	var outerUser User
	outerUserStructPtr = &outerUser

	err = json.Unmarshal(marshalledUserStruct, outerUserStructPtr)
	if err != nil {
		return nil, errors.New("could not Unmarshal User Struct")
	}

	// PUT THE ROOT KEY IN THE OUTER USER STRUCT
	outerUser.rootKey = rootKey

	// Check Integrity of encrypted INNER userstruct
	datastorePointerToInnerUserStruct := outerUser.EncInnerUserStructPointer

	encryptedInnerUserStruct, ok := userlib.DatastoreGet(datastorePointerToInnerUserStruct)
	if !ok {
		return nil, errors.New("issue with getting Inner User Struct in Datastore")
	}

	HMACEncryptedInnerUserStruct, err := userlib.HMACEval(innerUserStructHMACKey, encryptedInnerUserStruct)
	if err != nil {
		return nil, err
	}

	isHMACEqual := userlib.HMACEqual(outerUser.HMACEncInnerUserStruct, HMACEncryptedInnerUserStruct)

	// THIS IS OUR PASSWORD CHECK. THE HMAC KEY DEPENDS ON THE PASSWORD, SO IF
	// THE HMAC DOESN'T EQUAL THEN THE PASSWORD WAS WRONG
	// NOTE: THIS IS ALSO OUR TAMPERING CHECK
	if !isHMACEqual {
		return nil, errors.New("either the password is wrong or struct is tampered alert")
	}

	return outerUserStructPtr, nil
}

/* Basically the same as LoadOuterUserStructFromDatastore. Will error if HMAC is invalid */
func LoadInnerUserStructFromDatastore(rootKey []byte) (innerUserStructPtr *InnerUser, err error) {

	// Try loading the outer user struct, to check integrity of Inner User struct
	_, err = LoadOuterUserStructFromDatastore(rootKey)
	if err != nil {
		return nil, err
	}

	// Get the keys we need for the InnerUser Struct
	_, innerUserStructDecryptKey, _, err := GenerateUserStructKeys(rootKey)
	if err != nil {
		return nil, err
	}

	// Get the UUID's of the User Structs
	_, datastorePointerToInnerUserStruct, err := GetDatastorePointersToUserStructs(rootKey)
	if err != nil {
		return nil, err
	}

	// Get the Encrypted Inner User Struct from the Datastore
	encryptedInnerUserStruct, ok := userlib.DatastoreGet(datastorePointerToInnerUserStruct)
	if !ok {
		return nil, errors.New("could not get Encrypted Inner User Struct from Datastore")
	}

	// Decrypt and Unmarshal the Inner User Struct
	marshalledInnerUserStruct := userlib.SymDec(innerUserStructDecryptKey, encryptedInnerUserStruct)

	var innerUserStruct InnerUser
	innerUserStructPtr = &innerUserStruct

	err = json.Unmarshal(marshalledInnerUserStruct, innerUserStructPtr)
	if err != nil {
		return nil, errors.New("could not Unmarshal the Inner User Struct")
	}

	// decryptedInnerUserStructI, err := DecryptUnmarshalInnerStruct(encryptedInnerUserStruct, innerUserStructDecryptKey)
	// if err != nil {
	// 	return nil, err
	// }
	// innerUser, ok := decryptedInnerUserStructI.(InnerUser)
	// if !ok {
	// 	return nil, errors.New("cannot Type Cast the Struct to Inner User")
	// }

	// Return
	return innerUserStructPtr, nil
}

/* Returns err if failed at any point (including the HMAC check). Returns the pointer to the outer file struct */
func LoadOuterFileStructFromDatastore(filename string, fileRootKey []byte, revokeNumber int, linkedListNumber int) (outerFileStructPtr *OuterFile, err error) {

	// STEP 1: Generate Keys
	_, _, innerFileStructHMACKey, err := GenerateFileStructKeys(fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}

	// STEP 2: Generate Datastore UUID's (locations) for Outer + Inner File struct
	datastorePointerToOuterFileStruct, datastorePointerToInnerFileStruct, err := GetDatastorePointersToFileStructs(fileRootKey, filename, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}

	//userlib.DebugMsg("%v", datastorePointerToOuterFileStruct)

	// STEP 3: Load and Unmarshal the OuterFile struct from the Datastore
	marshalledOuterFileStruct, ok := userlib.DatastoreGet(datastorePointerToOuterFileStruct)
	if !ok {
		return nil, errors.New("could not Load the Outer File struct from the Datastore")
	}

	var outerFileStruct OuterFile
	outerFileStructPtr = &outerFileStruct

	err = json.Unmarshal(marshalledOuterFileStruct, outerFileStructPtr)
	if err != nil {
		return nil, errors.New("could not Unmarshal the Outer File struct from the Datastore")
	}

	// STEP 4: Load in the (Encrypted) InnerFile struct from
	// the Datastore using the pointer
	encryptedInnerFileStruct, ok := userlib.DatastoreGet(datastorePointerToInnerFileStruct)
	if !ok {
		return nil, errors.New("could not Load the Encrypted InnerFile struct from the Datastore")
	}

	// STEP 5: Compute the HMAC of the (Encrypted) InnerFile struct
	innerFileStructHMAC, err := userlib.HMACEval(innerFileStructHMACKey, encryptedInnerFileStruct)
	if err != nil {
		return nil, errors.New("could not compute the HMAC of the Inner File Struct ")
	}

	// STEP 6: Check the HMAC is correct for the InnerFile struct
	isEqual := userlib.HMACEqual(outerFileStruct.HMACEncInnerFileStruct, innerFileStructHMAC)
	if !isEqual {
		return nil, errors.New("the InnerFile struct has been TAMPERED with or file does not exist")
	}

	return outerFileStructPtr, nil
}

/* Basically the same as LoadOuterFileStructFromDatastore. Will error if HMAC is invalid. */
func LoadInnerFileStructFromDatastore(filename string, fileRootKey []byte, revokeNumber int, linkedListNumber int) (innerFileStructPtr *InnerFile, err error) {

	// First, check integrity of the InnerFile struct by trying to load the OuterFile struct and checking the HMAC
	_, err = LoadOuterFileStructFromDatastore(filename, fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}

	// Get the relevant keys
	_, innerFileStructDecryptKey, _, err := GenerateFileStructKeys(fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}

	// Get the relevant pointers
	_, datastorePointerToInnerFileStruct, err := GetDatastorePointersToFileStructs(fileRootKey, filename, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}

	// Now that integrity was checked, proceed to load the Encrypted InnerFile Struct
	encryptedInnerFileStruct, ok := userlib.DatastoreGet(datastorePointerToInnerFileStruct)
	if !ok {
		return nil, errors.New("could not Load the Encrypted InnerFile struct from the Datastore")
	}

	// Need to decrypt the inner struct
	marshalledInnerStruct := userlib.SymDec(innerFileStructDecryptKey, encryptedInnerFileStruct)

	// Then, need to unmarshal the inner struct
	var innerFileStruct InnerFile
	innerFileStructPtr = &innerFileStruct

	err = json.Unmarshal(marshalledInnerStruct, innerFileStructPtr)
	if err != nil {
		return nil, errors.New("could not Unmarshal the Inner Struct")
	}

	return innerFileStructPtr, nil
}

/* Key input may be any length greater than or equal to 16 */
func GenerateUUID(key []byte) (location UUID, err error) {
	// UUID = uuid.FromBytes(HASH(K)), [where K = HashKDF(SourceKey, purpose)]
	// SourceKey maybe a root key or a sub-root key

	location, err = uuid.FromBytes(userlib.Hash(key)[:16])

	if err != nil {
		return uuid.Nil, errors.New("cannot Generate UUID")
	}

	return location, err
}

/* Generate the keys we need for the (Outer) User Struct and the Inner User Struct */
func GenerateUserStructKeys(rootKey []byte) (outerUserStructKey []byte, innerUserStructEncryptDecryptKey []byte, innerUserStructHMACKey []byte, err error) {

	// Make the outerUserStructKey
	outerUserStructKey, err = userlib.HashKDF(rootKey, []byte("Key of Outer User Struct"))
	if err != nil {
		return nil, nil, nil, errors.New("could not Generate Key for Outer User Struct")
	}
	outerUserStructKey = outerUserStructKey[:16] // need to get first 16 bytes

	// Make the innerUserStructEncryptDecryptKey
	innerUserStructEncryptDecrypt64, err := userlib.HashKDF(outerUserStructKey, []byte("Key to Encrypt and Decrypt Inner User Struct"))
	if err != nil {
		return nil, nil, nil, errors.New("could not Generate Encrypt and Decrypt Key for Inner User Struct")
	}
	innerUserStructEncryptDecryptKey = innerUserStructEncryptDecrypt64[:16]

	// Make the innerUserStructHMACKey
	innerUserStructHMAC64, err := userlib.HashKDF(outerUserStructKey, []byte("Key to HMAC Encrypted Inner User Struct"))
	if err != nil {
		return nil, nil, nil, errors.New("could not Generate HMAC Key for Encrypted Inner User Struct")
	}
	innerUserStructHMACKey = innerUserStructHMAC64[:16] // get the first 16 bytes.

	return outerUserStructKey, innerUserStructEncryptDecryptKey, innerUserStructHMACKey, nil
}

func GetDatastorePointersToUserStructs(rootKey []byte) (datastorePointerToOuterUserStruct UUID, datastorePointerToInnerUserStruct UUID, err error) {

	// Get the keys we need for the InnerUser Struct
	outerUserStructKey, innerUserStructEncryptKey, _, err := GenerateUserStructKeys(rootKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Get the UUID of the (Outer) User struct
	datastorePointerToOuterUserStruct, err = GenerateUUID(outerUserStructKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("could not generate UUID for Outer user Struct")
	}

	// Get the UUID of the InnerUser struct
	datastorePointerToInnerUserStruct, err = GenerateUUID(innerUserStructEncryptKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("could not generate UUID for Inner user Struct")
	}

	return datastorePointerToOuterUserStruct, datastorePointerToInnerUserStruct, nil
}

func GetDatastorePointersToFileStructs(fileRootKey []byte, filename string, revokeNumber int, linkedListNumber int) (datastorePointerToOuterFileStruct UUID, datastorePointerToInnerFileStruct UUID, err error) {

	// STEP 1: Generate the relevant File Struct keys
	outerFileStructKey, innerFileStructEncryptKey, _, err := GenerateFileStructKeys(fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// STEP 2: Generate Datastore UUID's (locations) for Outer + Inner File struct
	datastorePointerToOuterFileStruct, err = GenerateUUID(outerFileStructKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("generate UUID failed for Outer File Struct")
	}
	datastorePointerToInnerFileStruct, err = GenerateUUID(innerFileStructEncryptKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("generate UUID failed for Inner File Struct")
	}

	return datastorePointerToOuterFileStruct, datastorePointerToInnerFileStruct, nil
}

/* For some user, say Alice, compute Alice's userFileRootKey (the fileRootKey in her personal namespace) */
func ComputeUserFileRootKey(rootKey []byte, filename string) (userFileRootKey []byte, err error) {

	//Important: the UserFileRootKey depends on the Revoke Number
	revokeNumber, err := GetRevokeNumberFromDatastore(rootKey, filename)
	if err != nil {
		return nil, err
	}

	purpose := "Compute userFileRootKey for file: " + filename + strconv.Itoa(revokeNumber)
	userFileRootKey, err = userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return nil, errors.New("could not compute User File Root Key")
	}
	userFileRootKey = userFileRootKey[:16] // get first 16 bytes

	return userFileRootKey, nil
}

/*
Generate the keys we need for the Outer File Struct and Inner File Struct.
revokeNumber is a number >= 0
linkedListNumber: {0 = sentinel, 1 = initial outer file struct, 2 = appended outer file struct ... }
*/
func GenerateFileStructKeys(fileRootKey []byte, revokeNumber int, linkedListNumber int) (outerFileStructKey []byte, innerFileStructEncryptDecryptKey []byte, innerFileStructHMACKey []byte, err error) {

	// Make the outerFileStructKey
	revokeNumberString := strconv.Itoa(revokeNumber)
	linkedListNumberString := strconv.Itoa(linkedListNumber)
	purpose := []byte("Key for Outer File Struct: " + revokeNumberString + linkedListNumberString)
	outerFileStructKey, err = userlib.HashKDF(fileRootKey, purpose)
	if err != nil {
		return nil, nil, nil, errors.New("could not generate Outer File Struct Key")
	}
	outerFileStructKey = outerFileStructKey[:16]

	// Make the innerFileStructEncryptDecryptKey
	purpose = []byte("Generate Key to Encrypt and Decrypt Inner File Struct")
	innerFileStructEncryptDecryptKey, err = userlib.HashKDF(outerFileStructKey, purpose)
	if err != nil {
		return nil, nil, nil, errors.New("could not generate Inner File Struct encrypt and decrypt Key")
	}
	innerFileStructEncryptDecryptKey = innerFileStructEncryptDecryptKey[:16]

	// Make the innerFileStructHMACKey
	purpose = []byte("Generate HMAC Key for Inner File Struct")
	innerFileStructHMACKey, err = userlib.HashKDF(outerFileStructKey, purpose)
	if err != nil {
		return nil, nil, nil, errors.New("could not generate Inner File Struct HMAC Key")
	}
	innerFileStructHMACKey = innerFileStructHMACKey[:16]

	return outerFileStructKey, innerFileStructEncryptDecryptKey, innerFileStructHMACKey, nil
}

/////////////////////////////////////
//  CREATING INVITATION FUNCTIONS  //
/////////////////////////////////////

// /* Function that returns true if the person with rootKey is the owner of the file. Owner has nothing stored in the location given by the GenerateUUID(userFileRootKey) */
// func isOwner(rootKey []byte, filename string) (isOwner bool, err error) {

// 	// Get the UserFileRootKey
// 	userFileRootKey, err := ComputeUserFileRootKey(rootKey, filename)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Get the datastore pointer to the pointer to the encryptedSharedFileRootKey
// 	datastorePointerToEncPointerToEncSharedFileRootKey, err := GenerateUUID(userFileRootKey)
// 	if err != nil {
// 		return false, errors.New("could not generate Datastore pointer to Shared File Root Key Pointer")
// 	}

// 	// Get the (Encrypted) SharedFileRootKeyPointer
// 	_, ok := userlib.DatastoreGet(datastorePointerToEncPointerToEncSharedFileRootKey)
// 	return !ok, nil
// }

/* Store the hashed recipientUsername in the Datastore in the array */
func StoreUsernameOfRecipientToDatastore(rootKey []byte, filename string, recipientUsername string) (err error) {

	// Get the location of the recipient's Username stored in the Datastore (this depends on the filename)
	purpose := "store username of recipient to Datastore " + filename
	keyToGenerateLocationOfRecipientsStruct, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return errors.New("could not generate the key for location of Recipient Username")
	}

	datastorePointerToRecipientsStruct, err := GenerateUUID(keyToGenerateLocationOfRecipientsStruct)
	if err != nil {
		return errors.New("could not generate the UUID for the location of Recipient Username")
	}

	// // Get the encrypt and decrypt Key for the array
	// recipientsStructEncryptDecryptKey, err = GenerateRecipientUsernamesEncryptDecryptKey(rootKey, filename)
	// if err != nil {
	// 	return err
	// }

	var recipientArray [][]byte

	// Load the Array and then append to it
	marshalledRecipientsStruct, ok := userlib.DatastoreGet(datastorePointerToRecipientsStruct)

	if !ok {

		// If the struct is not there, need to make it first

		recipientArray = make([][]byte, 0)
		recipientArray = append(recipientArray, userlib.Hash([]byte(recipientUsername)))

		recipientsStruct := RecipientUsernames{
			ListOfHashedRecipientUsernames: recipientArray,
		}

		marshalledRecipientsStruct, err := json.Marshal(recipientsStruct)
		if err != nil {
			return errors.New("could not marshal the Recipients Struct")
		}

		userlib.DatastoreSet(datastorePointerToRecipientsStruct, marshalledRecipientsStruct)

		//recipientArrayAsBytes :=
		//encryptedRecipientsArray = userlib.SymEnc(recipientsArrayEncryptDecryptKey, userlib.RandomBytes(16), recipientArray)

	} else {

		// If the struct is there, need to unmarshal, modify, then re-store it

		// Unmarshal
		var recipientsStruct RecipientUsernames
		err = json.Unmarshal(marshalledRecipientsStruct, &recipientsStruct)
		if err != nil {
			return errors.New("could not unmarshal the recipientsStruct from the Datastore")
		}

		// Modify
		recipientArray := recipientsStruct.ListOfHashedRecipientUsernames
		recipientArray = append(recipientArray, userlib.Hash([]byte(recipientUsername)))

		// Store
		recipientsStruct = RecipientUsernames{
			ListOfHashedRecipientUsernames: recipientArray,
		}

		marshalledRecipientsStruct, err := json.Marshal(recipientsStruct)
		if err != nil {
			return errors.New("could not marshal the Recipients Struct")
		}

		userlib.DatastoreSet(datastorePointerToRecipientsStruct, marshalledRecipientsStruct)

	}

	return nil
}

func LoadRecipientUsernamesStructFromDatastore(rootKey []byte, filename string) (recipientsStructPtr *RecipientUsernames, err error) {

	// Get the location of the recipient's Username stored in the Datastore (this depends on the filename)
	purpose := "store username of recipient to Datastore " + filename
	keyToGenerateLocationOfRecipientsStruct, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return nil, errors.New("could not generate the key for location of Recipient Username")
	}

	datastorePointerToRecipientsStruct, err := GenerateUUID(keyToGenerateLocationOfRecipientsStruct)
	if err != nil {
		return nil, errors.New("could not generate the UUID for the location of Recipient Username")
	}

	// Load the recipientUsernames struct from the Datastore
	marshalledRecipientsStruct, ok := userlib.DatastoreGet(datastorePointerToRecipientsStruct)

	if !ok {
		return nil, errors.New("could not load the RecipientUsernames struct from the Datastore")
	}

	// Unmarshal
	var recipientsStruct RecipientUsernames
	recipientsStructPtr = &recipientsStruct
	err = json.Unmarshal(marshalledRecipientsStruct, recipientsStructPtr)
	if err != nil {
		return nil, errors.New("could not unmarshal the recipientsStruct from the Datastore")
	}

	return recipientsStructPtr, nil
}

/* When creating an invitation, Alice needs to generate locations of the inner and outer invitation structs based uniquely on the FILE and the SHARED USER */
func GetDatastorePointersToInvitationStructs(rootKey []byte, filename string, recipientUsername string) (datastorePointerToOuterInvitationStruct UUID, datastorePointerToInnerInvitationStruct UUID, err error) {

	// Generate the UUID for the OuterInvitation struct
	purpose := "make key for generating the UUID of the Outer Invitation Struct " + filename + recipientUsername
	keyToGenerateOuterInvitationStructUUID, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("hashKDF could not generate the key used to make Outer InnerInvitation Struct UUID")
	}

	datastorePointerToOuterInvitationStruct, err = GenerateUUID(keyToGenerateOuterInvitationStructUUID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("could not generate the UUID for the Outer Invitation Struct")
	}

	// Get the Datastore pointer to the encrypted InnerInvitation Struct
	purpose = "make key for generating the UUID of encrypted Inner Invitation Struct " + filename + recipientUsername
	keyToGenerateEncryptedInnerInvitationStructUUID, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("hashKDF could not generate the key used to make Encrypted InnerInvitation Struct UUID")
	}

	// Generate a UUID in DataStore of InnerInvitation Struct
	datastorePointerToInnerInvitationStruct, err = GenerateUUID(keyToGenerateEncryptedInnerInvitationStructUUID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("failed to generate UUID for encrypted Invitation Struct")
	}

	return datastorePointerToOuterInvitationStruct, datastorePointerToInnerInvitationStruct, nil
}

/* Create and store the invitation structs, at the locations given by GetDatastorePointersToInvitationStructs() */
func CreateAndStoreInnerAndOuterInvitationStructs(rootKey []byte, sharedFileRootKey []byte, filename string, recipientUsername string) (datastorePointerToInnerInvitationStruct UUID, datastorePointerToOuterInvitationStruct UUID, err error) {

	// Get the public keys for the recipient -> (PKEEncKey PKEEncKey, DSVerifyKey DSVerifyKey, err error)
	PKEEncKey, _, err := GetUserPublicKeysFromUsername(recipientUsername)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Get the private keys for the SENDER (from the SENDER rootKey)
	// The SENDER needs to sign the innerInvitation struct
	_, DSSignKey, err := GetUserPrivateKeys(rootKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	datastorePointerToOuterInvitationStruct, datastorePointerToInnerInvitationStruct, err = GetDatastorePointersToInvitationStructs(rootKey, filename, recipientUsername)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// IMPORTANT: YOU CAN ASSUME THAT THE datastorePointerToEncSharedFileRootKey IS THE SAME ONE RECEIVED UPON INVITATION (OR A NEW ONE) AND IT IS THEREFORE OKAY TO OVERWRITE WITH THE (SAME) SHAREDFILEROOTKEY
	datastorePointerToEncSharedFileRootKey, err := GenerateOrGetDatastorePointerToEncSharedFileRootKey(rootKey, filename, userlib.Hash([]byte(recipientUsername)))
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Now, need to store the EncSharedFileRootKey at the location datastorePointerToEncSharedFileRootKey. First need to encrypt it
	encSharedFileRootKey := sharedFileRootKey
	//encSharedFileRootKey, err := userlib.PKEEnc(PKEEncKey, sharedFileRootKey)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("could not encrypt the sharedFileRootKey")
	}
	userlib.DatastoreSet(datastorePointerToEncSharedFileRootKey, encSharedFileRootKey)

	// Make the innerInvitationStruct
	innerInvitationStruct := InnerInvitation{
		PointerToEncSharedFileRootKey: datastorePointerToEncSharedFileRootKey,
	}

	// Marshal, Encrypt, and store InnerInvitation Struct in Datastore
	_, encryptedInnerInvitationStruct, err := MarshalEncryptStoreInnerStruct(innerInvitationStruct, PKEEncKey, datastorePointerToInnerInvitationStruct)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Sign the ENCRYPTED innerInvitationStruct
	signatureOfEncryptedInnerInvitationStruct, err := userlib.DSSign(DSSignKey, encryptedInnerInvitationStruct)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("could not sign the encrypted Inner Invitation Struct")
	}

	// Create OuterInvitationStruct
	outerInvitationStruct := OuterInvitation{
		EncInnerInvitationStructPointer: datastorePointerToInnerInvitationStruct,
		DigitalSignature:                signatureOfEncryptedInnerInvitationStruct,
	}

	_, err = MarshalStoreOuterStruct(outerInvitationStruct, datastorePointerToOuterInvitationStruct)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return datastorePointerToInnerInvitationStruct, datastorePointerToOuterInvitationStruct, nil
}

// /* UTILITY FUNCTION: Only use this function when creating an invitation or revoking. It generates the datastorePointerToEncSharedFileRootKey which is put in the InnerInvitation Struct */
// func GenerateDatastorePointerToEncSharedFileRootKey(rootKey []byte, filename string, recipientUsername string) (datastorePointerToEncSharedFileRootKey UUID, err error) {

// 	purpose := "make key for generating the UUID of Encrypted SharedFileRootKey " + filename + recipientUsername
// 	keyToGenerateSharedFileRootKeyUUID, err := userlib.HashKDF(rootKey, []byte(purpose))
// 	if err != nil {
// 		return uuid.Nil, errors.New("hashKDF could not generate the key used to make Encrypted SharedFileRootKey Pointer UUID")
// 	}

// 	datastorePointerToEncSharedFileRootKey, err = GenerateUUID(keyToGenerateSharedFileRootKeyUUID)
// 	if err != nil {
// 		return uuid.Nil, errors.New("could not generate UUID for Encrypted SharedFileRootKey")
// 	}

// 	return datastorePointerToEncSharedFileRootKey, nil
// }

/*
If a person is the OWNER of the file, they need to GENERATE a new pointer to the EncSharedFileRootKey. Otherwise, the person is NOT THE OWNER, so they need to GET the pointer they received upon invitation.
The recipientUsername is necessary in the case of generating a new pointer
*/
func GenerateOrGetDatastorePointerToEncSharedFileRootKey(rootKey []byte, filename string, hashedRecipientUsername []byte) (datastorePointerToEncSharedFileRootKey UUID, err error) {

	// Get the UserFileRootKey
	userFileRootKey, err := ComputeUserFileRootKey(rootKey, filename)
	if err != nil {
		return uuid.Nil, err
	}

	// Get the datastore pointer to the pointer to the encryptedSharedFileRootKey
	datastorePointerToEncPointerToEncSharedFileRootKey, err := GenerateUUID(userFileRootKey)
	if err != nil {
		return uuid.Nil, errors.New("could not generate Datastore pointer to Shared File Root Key Pointer")
	}

	// Get the (Encrypted) SharedFileRootKeyPointer
	encPointerToEncSharedFileRootKey, ok := userlib.DatastoreGet(datastorePointerToEncPointerToEncSharedFileRootKey)
	isOwner := !ok

	if isOwner {

		purpose := "make key for generating the UUID of Encrypted SharedFileRootKey " + filename + string(hashedRecipientUsername)
		keyToGenerateSharedFileRootKeyUUID, err := userlib.HashKDF(rootKey, []byte(purpose))
		if err != nil {
			return uuid.Nil, errors.New("hashKDF could not generate the key used to make Encrypted SharedFileRootKey Pointer UUID")
		}

		datastorePointerToEncSharedFileRootKey, err = GenerateUUID(keyToGenerateSharedFileRootKeyUUID)
		if err != nil {
			return uuid.Nil, errors.New("could not generate UUID for Encrypted SharedFileRootKey")
		}

	} else {

		// Get the (PRIVATE) decryption key for the user
		PKEDecKey, _, err := GetUserPrivateKeys(rootKey)
		if err != nil {
			return uuid.Nil, err
		}

		// Decrypt the encPointerToEncSharedFileRootKey to bytes
		bytePointerToEncSharedFileRootKey, err := userlib.PKEDec(PKEDecKey, encPointerToEncSharedFileRootKey)
		if err != nil {
			return uuid.Nil, errors.New("could not decrypt the encryptedSharedFileRootKeyPointer")
		}

		// convert the byte array to UUID
		// https://github.com/tidwall/gjson/issues/107
		// Note: to encrypt the uuid, do:
		//userlib.PKEEnc(PKEEncKey, []byte(datastorePointerToSharedFileRootKey.String()))
		datastorePointerToEncSharedFileRootKey, err = uuid.ParseBytes(bytePointerToEncSharedFileRootKey)
		if err != nil {
			return uuid.Nil, errors.New("could not convert from bytes to UUID using ParseBytes")
		}

	}

	return datastorePointerToEncSharedFileRootKey, nil
}

/////////////////////////////////////
// ACCEPTING INVITATION FUNCTIONS  //
/////////////////////////////////////

/* Only use this function when accepting an invitation. The receiver stores the content of the InnerInvitation Struct in his own "personal namespace" */
func StorePointerToEncSharedFileRootKey(rootKey []byte, filename string, pointerToEncSharedFileRootKey UUID) (err error) {

	// Get the UserFileRootKey
	userFileRootKey, err := ComputeUserFileRootKey(rootKey, filename)
	if err != nil {
		return err
	}

	// Get user's Public Enc Key
	PKEEncKey, _, err := GetUserPublicKeys(rootKey)
	if err != nil {
		return err
	}

	// Get the datastore pointer to the pointer to the encryptedSharedFileRootKey
	datastorePointerToEncPointerToEncSharedFileRootKey, err := GenerateUUID(userFileRootKey)
	if err != nil {
		return errors.New("could not generate Datastore pointer to Shared File Root Key Pointer")
	}

	// NEW POINTER STORING
	// Store the contents at the location, but first need to MARSHAL IT
	bytePointerToEncSharedFileRootKey, err := json.Marshal(pointerToEncSharedFileRootKey)
	if err != nil {
		return errors.New("could not marshal the pointerToEncSharedFileRootKey")
	}

	encPointerToEncSharedFileRootKey, err := userlib.PKEEnc(PKEEncKey, bytePointerToEncSharedFileRootKey)
	if err != nil {
		return errors.New("could not encrypt the encSharedFileRootKeyPointer")
	}

	userlib.DatastoreSet(datastorePointerToEncPointerToEncSharedFileRootKey, encPointerToEncSharedFileRootKey)

	return nil
}

/////////////////////////////////////
// REVOKING INVITATION FUNCTIONS   //
/////////////////////////////////////

func SetRevokeNumberInDatastore(rootKey []byte, filename string, revokeNumber int) (err error) {

	purpose := "get the location of the revoke number from the Datastore " + filename
	keyToGetLocationOfRevokeNumber, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return errors.New("could not generate the key for location of RevokeNumber")
	}

	datastorePointerToRevokeNumber, err := GenerateUUID(keyToGetLocationOfRevokeNumber)
	if err != nil {
		return errors.New("could not generate UUID of the RevokeNumber in Datastore")
	}

	revokeNumberBytes := []byte(strconv.Itoa(revokeNumber))

	userlib.DatastoreSet(datastorePointerToRevokeNumber, revokeNumberBytes)

	return
}

/* Get the revoke number from the datastore. If it does not exist in Datastore, set revokeNumber to 0 (revoke hasn't happened yet). If an error occurred elsewhere, return -1 and error*/
func GetRevokeNumberFromDatastore(rootKey []byte, filename string) (revokeNumber int, err error) {

	purpose := "get the location of the revoke number from the Datastore " + filename
	keyToGetLocationOfRevokeNumber, err := userlib.HashKDF(rootKey, []byte(purpose))
	if err != nil {
		return -1, errors.New("could not generate the key for location of RevokeNumber")
	}

	datastorePointerToRevokeNumber, err := GenerateUUID(keyToGetLocationOfRevokeNumber)
	if err != nil {
		return -1, errors.New("could not generate UUID of the RevokeNumber in Datastore")
	}

	// If the revokeNumber does not exist, then !ok is true, so return 0 and nil
	revokeNumberBytes, ok := userlib.DatastoreGet(datastorePointerToRevokeNumber)
	if !ok {
		return 0, nil
	}

	// If the revokeNumber does exist, then !ok is false, so return the number
	revokeNumber, err = strconv.Atoi(string(revokeNumberBytes))
	if err != nil {
		return -1, errors.New("could not convert revokeNumber from bytes to int")
	}

	return revokeNumber, nil
}

/*
Returns true if the file with filename, created by a user with the provided rootKey, is shared with the recipientUsername (the invitation struct exists and is not garbage).
If an error occurs, fileIsSharedWithUser is set to false by DEFAULT (but could be set to true)
*/
func FileIsSharedWithUser(rootKey []byte, filename string, recipientUsername string) (fileIsSharedWithUser bool, err error) {

	datastorePointerToOuterInvitationStruct, _, err := GetDatastorePointersToInvitationStructs(rootKey, filename, recipientUsername)
	if err != nil {
		return false, err
	}

	encryptedMarshalledOuterInvitationStruct, ok := userlib.DatastoreGet(datastorePointerToOuterInvitationStruct)
	if !ok {
		return false, nil
	}

	garbage := []byte("you got pranked")
	equal := bytes.Equal(encryptedMarshalledOuterInvitationStruct, garbage)

	return !equal, nil
}

/*
Only use this function after accepting an invitation (for example, in StoreFile, LoadFile, AppendToFile).
This function takes in the rootKey of the RECEIVER (Bob), and returns the pointer to sharedFileRootKey
If the sharedFileRootKey does not exist for the file, return the userFileRootKey.
The owner will not have his own fileRootKey stored in the Datastore.
*/
func GetSharedFileRootKey(rootKey []byte, filename string) (sharedFileRootKey []byte, err error) {

	// Get the UserFileRootKey
	userFileRootKey, err := ComputeUserFileRootKey(rootKey, filename)
	if err != nil {
		return nil, err
	}

	// Get the (PRIVATE) decryption key for the user
	PKEDecKey, _, err := GetUserPrivateKeys(rootKey)
	if err != nil {
		return nil, err
	}

	// Get the datastore pointer to the pointer to the encryptedSharedFileRootKey
	datastorePointerToEncPointerToEncSharedFileRootKey, err := GenerateUUID(userFileRootKey)
	if err != nil {
		return nil, errors.New("could not generate Datastore pointer to Shared File Root Key Pointer")
	}

	// Get the (Encrypted) SharedFileRootKeyPointer
	encPointerToEncSharedFileRootKey, ok := userlib.DatastoreGet(datastorePointerToEncPointerToEncSharedFileRootKey)
	//userlib.DebugMsg("Enc Pointer To Enc Shared File Root Key %v", encPointerToEncSharedFileRootKey)
	//return nil, errors.New("encryptedSharedFileRootKeyPointer does not exist in the Datastore")

	// Need to check two cases.
	// Case 1: The person has a file called 'filename' which has been shared with them (ok == true), and if so we should return the sharedFileRootKey
	// Case 2: The person does NOT have a file called 'filename' which has been shared with them (ok == false), in which case we should return the uesrFileRootKey and nil

	if ok {

		// Decrypt the encPointerToEncSharedFileRootKey to bytes
		bytePointerToEncSharedFileRootKey, err := userlib.PKEDec(PKEDecKey, encPointerToEncSharedFileRootKey)
		if err != nil {
			return nil, errors.New("could not decrypt the encryptedSharedFileRootKeyPointer")
		}

		// convert the byte array to UUID
		// https://github.com/tidwall/gjson/issues/107
		// Note: to encrypt the uuid, do:
		//userlib.PKEEnc(PKEEncKey, []byte(datastorePointerToSharedFileRootKey.String()))

		// NEW POINTER STORING
		var pointerToEncSharedFileRootKey UUID

		err = json.Unmarshal(bytePointerToEncSharedFileRootKey, &pointerToEncSharedFileRootKey)
		if err != nil {
			return nil, errors.New("could not unmarshal the bytePointerToEncSharedFileRootKey")
		}

		// pointerToEncSharedFileRootKey, err := uuid.ParseBytes(bytePointerToEncSharedFileRootKey)
		// if err != nil {
		// 	return nil, errors.New("could not convert from bytes to UUID using ParseBytes")
		// }

		// Get the (Encrypted) SharedFileRootKey
		sharedFileRootKey, ok = userlib.DatastoreGet(pointerToEncSharedFileRootKey)
		if !ok {
			return nil, errors.New("could not get the encryptedSharedFileRootKey from Datastore")
		}

		//userlib.DebugMsg("%v", PKEDecKey)

		// Decrypt the (Encrypted) SharedFileRootKey
		// sharedFileRootKey, err = userlib.PKEDec(PKEDecKey, encryptedSharedFileRootKey)
		// if err != nil {
		// 	return nil, errors.New("could not decrypt the sharedFileRootKey")
		// }

		//userlib.DebugMsg("SFRK: %v", sharedFileRootKey)

	} else {
		sharedFileRootKey = userFileRootKey
	}

	return sharedFileRootKey, nil
}

////////////////////////
// UTILITY FUNCTIONS  //
////////////////////////

// Function to check if a file exists in the User's personal namespace already. If it does, return true. Otherwise, return false. If an error occured, the exists bool is false by default.
func FileExistsInUsersPersonalNamespace(rootKey []byte, filename string) (fileExistsInUsersPersonalNamespace bool, err error) {

	fileRootKey, err := GetSharedFileRootKey(rootKey, filename)
	if err != nil {
		return false, err
	}

	// Get the relevant pointers for the SENTINEL (just try to see if the sentinel exists in the datastore)
	revokeNumber := 0
	linkedListNumber := 0
	_, datastorePointerToInnerFileStruct, err := GetDatastorePointersToFileStructs(fileRootKey, filename, revokeNumber, linkedListNumber)
	if err != nil {
		return false, err
	}

	// Try to get the Inner File Struct from the Datastore (could have also done the outer struct?)
	_, ok := userlib.DatastoreGet(datastorePointerToInnerFileStruct)

	// Return true if the Datatstore has the Inner File Struct
	return ok, nil
}

/* USE ONLY FOR THE INNER STRUCTS */
func MarshalEncryptStoreInnerStruct(innerStructI interface{}, innerStructEncryptKeyI interface{}, datastorePointerToInnerStruct UUID) (marshalledInnerStruct []byte, encryptedInnerStruct []byte, err error) {

	// Do any necessary type casts
	switch innerStructI.(type) {
	case User:
		return nil, nil, errors.New("must supply an INNER struct (actually supplied User)")
	case InnerUser:

		innerStruct, ok := innerStructI.(InnerUser)
		if !ok {
			return nil, nil, errors.New("the type cast to the correct struct failed")
		}
		// Type cast the encryption key to []byte
		innerStructEncryptKey, ok := innerStructEncryptKeyI.([]byte)
		if !ok {
			return nil, nil, errors.New("the encryption key should be of type []byte")
		}

		marshalledInnerStruct, err = json.Marshal(innerStruct)
		if err != nil {
			return nil, nil, errors.New("error in Marshal Struct")
		}
		encryptedInnerStruct = userlib.SymEnc(innerStructEncryptKey, userlib.RandomBytes(16), marshalledInnerStruct)
		userlib.DatastoreSet(datastorePointerToInnerStruct, encryptedInnerStruct)

		return marshalledInnerStruct, encryptedInnerStruct, nil
	case OuterFile:
		return nil, nil, errors.New("must supply an INNER struct (actually supplied OuterFile)")
	case InnerFile:

		innerStruct, ok := innerStructI.(InnerFile)
		if !ok {
			return nil, nil, errors.New("the type cast to the correct struct failed")
		}
		// Type cast the encryption key to []byte
		innerStructEncryptKey, ok := innerStructEncryptKeyI.([]byte)
		if !ok {
			return nil, nil, errors.New("the encryption key should be of type []byte")
		}

		marshalledInnerStruct, err = json.Marshal(innerStruct)
		if err != nil {
			return nil, nil, errors.New("error in Marshal Struct")
		}
		encryptedInnerStruct = userlib.SymEnc(innerStructEncryptKey, userlib.RandomBytes(16), marshalledInnerStruct)
		userlib.DatastoreSet(datastorePointerToInnerStruct, encryptedInnerStruct)

		return marshalledInnerStruct, encryptedInnerStruct, nil
	case OuterInvitation:
		return nil, nil, errors.New("must supply an INNER struct (actually supplied OuterInvitation)")
	case InnerInvitation:

		innerStruct, ok := innerStructI.(InnerInvitation)
		if !ok {
			return nil, nil, errors.New("the type cast to the correct struct failed")
		}
		// Type cast the encryption key to PKEDecKey
		innerStructEncryptKey, ok := innerStructEncryptKeyI.(PKEEncKey)
		if !ok {
			return nil, nil, errors.New("the encryption key should be of type PKEEncKey")
		}

		marshalledInnerStruct, err = json.Marshal(innerStruct)
		if err != nil {
			return nil, nil, errors.New("error in Marshal Struct")
		}
		encryptedInnerStruct, err = userlib.PKEEnc(innerStructEncryptKey, marshalledInnerStruct)
		if err != nil {
			return nil, nil, errors.New("could not encrypt the marshalled InnerInvitation Struct")
		}
		userlib.DatastoreSet(datastorePointerToInnerStruct, encryptedInnerStruct)

		return marshalledInnerStruct, encryptedInnerStruct, nil
	default:
		return nil, nil, errors.New("the argument was not a struct")
	}
}

/* USE ONLY FOR THE OUTER STRUCTS */
func MarshalStoreOuterStruct(outerStructI interface{}, datastorePointerToOuterStruct UUID) (marshalledOuterStruct []byte, err error) {

	// Do any necessary type casts
	switch outerStructI.(type) {
	case User:
		outerStruct, ok := outerStructI.(User)
		if !ok {
			return nil, errors.New("the type cast to the correct struct failed")
		}

		marshalledOuterStruct, err = json.Marshal(outerStruct)
		if err != nil {
			return nil, errors.New("error in Marshal Struct")
		}

		userlib.DatastoreSet(datastorePointerToOuterStruct, marshalledOuterStruct)

		return marshalledOuterStruct, nil
	case InnerUser:
		return nil, errors.New("must supply an OUTER struct (actually supplied InnerUser)")
	case OuterFile:
		outerStruct, ok := outerStructI.(OuterFile)
		if !ok {
			return nil, errors.New("the type cast to the correct struct failed")
		}

		marshalledOuterStruct, err = json.Marshal(outerStruct)
		if err != nil {
			return nil, errors.New("error in Marshal Struct")
		}

		userlib.DatastoreSet(datastorePointerToOuterStruct, marshalledOuterStruct)

		return marshalledOuterStruct, nil
	case InnerFile:
		return nil, errors.New("must supply an OUTER struct (actually supplied InnerFile)")
	case OuterInvitation:
		outerStruct, ok := outerStructI.(OuterInvitation)
		if !ok {
			return nil, errors.New("the type cast to the correct struct failed")
		}

		marshalledOuterStruct, err = json.Marshal(outerStruct)
		if err != nil {
			return nil, errors.New("error in Marshal Struct")
		}

		userlib.DatastoreSet(datastorePointerToOuterStruct, marshalledOuterStruct)

		return marshalledOuterStruct, nil
	case InnerInvitation:
		return nil, errors.New("must supply an OUTER struct (actually supplied InnerInvitation)")
	default:
		return nil, errors.New("the argument was not a struct")
	}
}

/* BUGGY FUNCTION???? DO NOT USE */
func DecryptUnmarshalInnerStruct(encryptedInnerStruct []byte, innerStructDecryptKey []byte) (plainTextInnerStruct interface{}, err error) {

	// First, need to decrypt the inner struct
	marshalledInnerStruct := userlib.SymDec(innerStructDecryptKey, encryptedInnerStruct)

	// Then, need to unmarshal the inner struct
	// var decryptedInnerStruct interface{} //don't know what type of struct it is yet
	innerStructPtr := &plainTextInnerStruct

	err = json.Unmarshal(marshalledInnerStruct, innerStructPtr)
	if err != nil {
		return nil, errors.New("could not Unmarshal the Inner Struct")
	}

	return plainTextInnerStruct, nil
}

////////////////////////////////////
// FUNCTIONS TO COMPLETE
////////////////////////////////////

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {

	// Step 1: Check if an empty username is provided. If so, error
	isEmpty := IsEmptyUsername(username)
	if isEmpty {
		return nil, errors.New("empty username was provided")
	}

	// Step 2: Check if a user with the given username exists.
	userExists := CheckUserExists(username, password)
	if userExists {
		return nil, errors.New("the user already exists")
	}

	// Step 3: Generate the Digital Signature and RSA Keys, and use them
	// to create the Inner and Outer User Structs
	_, userdataptr, err = CreateAndStoreInnerAndOuterUserStructs(username, password)

	return userdataptr, err
}

/* This function returns an (Outer) User Struct pointer, unmarshalled! */
func GetUser(username string, password string) (userdataptr *User, err error) {

	// Step 1: Check if a user with the given username exists. If doesn't exist, error
	userExists := CheckUserExists(username, password)
	if !userExists {
		return nil, errors.New("the user does not exist")
	}

	rootKey := userlib.Argon2Key([]byte(password), []byte(username), 16) //PBKDF

	// Step 2: Load the Outer User Struct from the Datastore, which includes integrity check
	userdataptr, err = LoadOuterUserStructFromDatastore(rootKey)
	if err != nil {
		return nil, err
	}

	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	// Get the FileRootKey
	rootKey := (*userdata).rootKey

	// Need to check two cases.
	// Case 1: The person has a file called 'filename' which has been shared with them, in which case GetSharedFileRootKey() will not error
	// Case 2: The person does NOT have a file called 'filename' which has been shared with them, in which case GetSharedFileRootKey() will error
	fileRootKey, err := GetSharedFileRootKey(rootKey, filename)
	if err != nil {
		return err
	}

	// STORE THE FIRST FILE STRUCTS
	revokeNumber := 0
	linkedListNumber := 1 // linkedListNumber MUST BE EUQAL TO 1 FOR FIRST FILE
	linkedListSize := 1   // since we are storing the first file structs, set size = 1
	err = CreateAndStoreInnerAndOuterFileStructs(filename, content, fileRootKey, revokeNumber, linkedListNumber, linkedListSize)
	if err != nil {
		return err
	}

	// STORE THE SENTINEL STRUCTS
	revokeNumber = 0
	linkedListNumber = 0 // linkedListNumber MUST BE EUQAL TO 0 FOR SENTINEL
	linkedListSize = 1   // since we already added the first file structs, set size = 1
	err = CreateAndStoreInnerAndOuterFileStructs(filename, content, fileRootKey, revokeNumber, linkedListNumber, linkedListSize)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) error {

	rootKey := (*userdata).rootKey

	fileRootKey, err := GetSharedFileRootKey(rootKey, filename)
	if err != nil {
		return err
	}

	// STEP 1: Load the Sentinel Inner File Struct from the Datastore. For the sentinel, linkedListNumber = 0
	revokeNumber := 0
	linkedListNumber := 0
	sentinelInnerFileStructPtr, err := LoadInnerFileStructFromDatastore(filename, fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return err
	}
	sentinelInnerFileStruct := (*sentinelInnerFileStructPtr)

	// STEP 2: Compute the UPDATED Linked List Size from the Sentinel
	linkedListSize := sentinelInnerFileStruct.LinkedListNumber + 1

	// STEP 3: Store the Sentinel Back into the Datastore, with the new linkedListSize.
	revokeNumber = 0
	linkedListNumber = 0
	err = CreateAndStoreInnerAndOuterFileStructs(filename, nil, fileRootKey, revokeNumber, linkedListNumber, linkedListSize)
	if err != nil {
		return err
	}

	// STEP 4: Create the new Outer and Inner File Structs, Storing them into the Datastore. Now, linkedListNumber = linkedListSize
	revokeNumber = 0
	linkedListNumber = linkedListSize
	err = CreateAndStoreInnerAndOuterFileStructs(filename, content, fileRootKey, revokeNumber, linkedListNumber, linkedListSize)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	rootKey := (*userdata).rootKey

	fileRootKey, err := GetSharedFileRootKey(rootKey, filename)
	if err != nil {
		return nil, err
	}

	// STEP 1: Load the Sentinel Inner File Struct from the Datastore. For the sentinel, linkedListNumber = 0
	revokeNumber := 0
	linkedListNumber := 0
	sentinelInnerFileStructPtr, err := LoadInnerFileStructFromDatastore(filename, fileRootKey, revokeNumber, linkedListNumber)
	if err != nil {
		return nil, err
	}
	sentinelInnerFileStruct := (*sentinelInnerFileStructPtr)

	// STEP 2: Get the LinkedListSize from the Sentinel
	linkedListSize := sentinelInnerFileStruct.LinkedListNumber

	// STEP 3: Loop over the stored Inner File Structs, loading the content from them.
	// The content is appended to a variable, called "contents"

	var contents []byte

	var i int
	for i = 1; i <= linkedListSize; i++ {

		revokeNumber = 0
		linkedListNumber = i
		innerFileStructPtr, err := LoadInnerFileStructFromDatastore(filename, fileRootKey, revokeNumber, linkedListNumber)
		if err != nil {
			return nil, err
		}
		innerFileStruct := (*innerFileStructPtr)

		// Append Contents
		contents = append(contents, innerFileStruct.Content...)
	}

	return contents, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {

	// The owner's rootKey
	rootKey := (*userdata).rootKey

	//Check the file integrity by trying to load the file
	_, err = userdata.LoadFile(filename)
	if err != nil {
		return uuid.Nil, errors.New("tampering of file has been detected")
	}

	// get the SHARED fileRootKey (Why is this? It is because a shared user can also invite people)
	// You can assume that CreateInvitation is called on a recipient who is already currently authorized to access the file (in other words, Shared File Root Key will always exist)
	sharedFileRootKey, err := GetSharedFileRootKey(rootKey, filename)
	if err != nil {
		return uuid.Nil, errors.New("failed to generate a fileRootKey")
	}

	// //sharedFileRootKey := fileRootKey // this is true because the owner of the file is sharing their own fileRootKey with another person

	_, datastorePointerToOuterInvitationStruct, err := CreateAndStoreInnerAndOuterInvitationStructs(rootKey, sharedFileRootKey, filename, recipientUsername)
	if err != nil {
		return uuid.Nil, err
	}

	err = StoreUsernameOfRecipientToDatastore(rootKey, filename, recipientUsername)
	if err != nil {
		return uuid.Nil, err
	}

	return datastorePointerToOuterInvitationStruct, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) (err error) {

	rootKey := (*userdata).rootKey

	// Check if Bob has the same name file in his personal namespace as sender
	fileExistsInUsersPersonalNamespace, err := FileExistsInUsersPersonalNamespace(rootKey, filename)

	// First, check if there was an error, and if so just return the error
	if err != nil {
		return errors.New("there is the same name file existing in receiver's namespace")
	}

	// If there was no error, now check if the file actually exists in the user's namespace
	if fileExistsInUsersPersonalNamespace {
		return errors.New("cannot accept invitation when file already exists in personal namespace")
	}

	// Step 1: Bob accesses the serialized struct using DatastoreGet by invitationPtr
	marshalledOuterInvitaionStruct, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("failed to get Outer Struct using invitation pointer")
	}

	var outerInvitationStruct OuterInvitation
	outerInvitationStructPtr := &outerInvitationStruct

	// Step 2: deserialized
	err = json.Unmarshal(marshalledOuterInvitaionStruct, outerInvitationStructPtr)
	if err != nil {
		return errors.New("failed to unmarshall Outer Invitation Struct")
	}

	// Access InnerInvitationStruct using OuterInvitationStruct variables
	encryptedInnerInvitationStructPointer := outerInvitationStruct.EncInnerInvitationStructPointer

	// Load the encrypted InnerInvitation Struct from the Datastore
	encryptedInnerInvitationStruct, ok := userlib.DatastoreGet(encryptedInnerInvitationStructPointer)
	if !ok {
		return errors.New("could not get the Encrypted InnerInvitation Struct from the Datastore")
	}

	// Before decrypting, check the signature, firstly get public DS verify key
	_, DSVerifyKey, err := GetUserPublicKeysFromUsername(senderUsername)
	if err != nil {
		return err
	}

	signature := outerInvitationStruct.DigitalSignature
	err = userlib.DSVerify(DSVerifyKey, encryptedInnerInvitationStruct, signature)
	if err != nil {
		return errors.New("invitation is tampered")
	}

	// Decrypt the InnerInvitation struct, get recipient's private decryption key
	PKEDecKey, _, err := GetUserPrivateKeys(rootKey)
	if err != nil {
		return errors.New("failed to get the recipient's private Decryption key")
	}

	// Decrypt and unMarshall the InnerInvitation Struct using recipient's private key; HAVE a decryptUnmarshallInnerStruct function
	marshalledInnerInvitationStruct, err := userlib.PKEDec(PKEDecKey, encryptedInnerInvitationStruct)
	if err != nil {
		return errors.New("failed to decrypt Inner Invitation Struct")
	}

	var innerInvitationStruct InnerInvitation
	innerInvitationStructPointer := &innerInvitationStruct
	err = json.Unmarshal(marshalledInnerInvitationStruct, innerInvitationStructPointer)
	if err != nil {
		return errors.New("failed to unmarshal the Inner Invitation Struct")
	}

	// Get the sharedFileRootKey from the innerInvitationStruct and put it in the Datastore in the correct location
	pointerToEncSharedFileRootKey := innerInvitationStruct.PointerToEncSharedFileRootKey

	err = StorePointerToEncSharedFileRootKey(rootKey, filename, pointerToEncSharedFileRootKey)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) (err error) {

	rootKey := (*userdata).rootKey

	garbage := []byte("you got pranked")

	// CHECK 1: The given filename does not exist in the caller’s personal file namespace = return error
	fileExistsInUsersPersonalNamespace, err := FileExistsInUsersPersonalNamespace(rootKey, filename)
	if err != nil {
		return errors.New("errors happen when checking files in user personal name space")
	}

	if !fileExistsInUsersPersonalNamespace {
		return errors.New("cannot revoke when you don't even have a file with that name")
	}

	// CHECK 2: The given filename is not currently shared with recipientUsername.
	fileIsSharedWithUser, err := FileIsSharedWithUser(rootKey, filename, recipientUsername)
	if err != nil {
		return errors.New("error happens when try to see if the file is shared with another user")
	}

	if !fileIsSharedWithUser {
		return errors.New("cannot revoke when you have not even shared the file with the other person")
	}

	// Step 3: Get locations of the invitation structs AND the sharedFileRootKey

	datastorePointerToOuterInvitationStruct, datastorePointerToInnerInvitationStruct, err := GetDatastorePointersToInvitationStructs(rootKey, filename, recipientUsername)
	if err != nil {
		return err
	}

	datastorePointerToEncSharedFileRootKey, err := GenerateOrGetDatastorePointerToEncSharedFileRootKey(rootKey, filename, userlib.Hash([]byte(recipientUsername)))
	if err != nil {
		return err
	}

	// Set the invitation to garbage
	userlib.DatastoreSet(datastorePointerToOuterInvitationStruct, garbage)
	userlib.DatastoreSet(datastorePointerToInnerInvitationStruct, garbage)

	// Set the Shared File Root Key ITSELF to garbage (it is stored at the location given by the shared pointer)
	// This change to the key should cascade down the invitation tree because only pointers to they key were shared and stored
	userlib.DatastoreSet(datastorePointerToEncSharedFileRootKey, garbage)

	///////////////////////////////////////////////
	// COPY THE OLD CONTENTS TO THE NEW LOCATION //
	///////////////////////////////////////////////

	content, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	// Change the REVOKE NUMBER. This is extremely important to do because it updates the location where the OWNER will store the new copy of the file
	// -- ComputeUserFileRootKey() function relies on the updated Revoke Number. This function is used whenever we store the file to the Datastore, for example.
	revokeNumber, err := GetRevokeNumberFromDatastore(rootKey, filename)
	if err != nil {
		return err
	}
	SetRevokeNumberInDatastore(rootKey, filename, revokeNumber+1)

	// Now, put a copy in the file in the Datastore by calling StoreFile.... (this is changing the location for the OWNER)
	err = userdata.StoreFile(filename, content)
	if err != nil {
		return err
	}

	///////////////////////////////////////////////////////////////////
	// CHANGE THE SHARED FILE ROOT KEY OF THE FILE FOR EVERYONE ELSE //
	///////////////////////////////////////////////////////////////////

	// Change the location for everyone ELSE who we have shared with, by getting their datastorePointerToEncSharedFileRootKey and changing it to the NEW value
	recipientsStructPtr, err := LoadRecipientUsernamesStructFromDatastore(rootKey, filename)
	if err != nil {
		return err
	}

	listOfHashedRecipientUsernames := (*recipientsStructPtr).ListOfHashedRecipientUsernames

	// Now, iterate through all of the shared people's usernames, get the DatastorePointerToEncSharedFileRootKey, and change it to the new one
	for _, hashedRecipientUsername := range listOfHashedRecipientUsernames {

		datastorePointerToEncSharedFileRootKey, err = GenerateOrGetDatastorePointerToEncSharedFileRootKey(rootKey, filename, userlib.Hash([]byte(hashedRecipientUsername)))
		if err != nil {
			return err
		}

		newSharedFileRootKey, err := ComputeUserFileRootKey(rootKey, filename)
		if err != nil {
			return err
		}

		userlib.DatastoreSet(datastorePointerToEncSharedFileRootKey, newSharedFileRootKey)

	}

	return nil
}
