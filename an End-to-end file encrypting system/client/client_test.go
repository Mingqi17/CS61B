package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	var eve *client.User
	var frank *client.User
	var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	dorisFile := "dorisFile.txt"
	eveFile := "eveFile.txt"
	frankFile := "frankFile.txt"
	graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser then GetUser", func() {

			userlib.DebugMsg(("Testing InitUser then GetUser"))
			_, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg(("Testing User Does Not Exist"))
			_, err = client.GetUser("bob", defaultPassword)
			Expect(err).NotTo(BeNil())

		})

		Specify("Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

		})

		Specify("Basic Test: Testing Single User Store/Load.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			//
			// _, err = bob.LoadFile(bobFile)
			// Expect(err).To(BeNil())
			//

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			// _, err = bob.LoadFile(bobFile)
			// Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Complicated Test: Testing Revoke Functionality Twice", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			// MAKE BOB, CHARLES, DORIS
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			/////////

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Charles creating invite for Doris for file %s, and Doris accepting invite under name %s.", charlesFile, dorisFile)
			invite, err = charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			//
			//_, err = charles.LoadFile(charlesFile)
			//Expect(err).To(BeNil())
			//

			err = doris.AcceptInvitation("charles", invite, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Doris can load the file.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			// MAKE EVE
			eve, err = client.InitUser("eve", defaultPassword)
			Expect(err).To(BeNil())
			/////////

			// Make frank and grace
			frank, err = client.InitUser("frank", defaultPassword)
			Expect(err).To(BeNil())

			grace, err = client.InitUser("grace", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Eve for file %s, and Eve accepting invite under name %s.", aliceFile, eveFile)

			invite, err = alice.CreateInvitation(aliceFile, "eve")
			Expect(err).To(BeNil())

			err = eve.AcceptInvitation("alice", invite, eveFile)
			Expect(err).To(BeNil())

			// Invite frank, grace, they accept

			userlib.DebugMsg("Alice creating invite for Frank and Grace for file %s, and they accepting invite under names.", aliceFile)

			invite, err = alice.CreateInvitation(aliceFile, "frank")
			Expect(err).To(BeNil())

			err = frank.AcceptInvitation("alice", invite, frankFile)
			Expect(err).To(BeNil())

			invite, err = alice.CreateInvitation(aliceFile, "grace")
			Expect(err).To(BeNil())

			err = grace.AcceptInvitation("alice", invite, graceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Eve appending to file %s, content: %s", eveFile, contentTwo)
			err = eve.AppendToFile(eveFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			// Make aliceDesktop now
			userlib.DebugMsg("Getting user Alice (aliceDesktop)")
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that eve sees expected file data.")
			data, err = eve.LoadFile(eveFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that doris sees expected file data.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("bob appending to file %s, content: %s", bobFile, contentThree)
			err = bob.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())

			//REVOKE BOB
			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentThree)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that Doris cannot store to the file.")
			err = doris.StoreFile(dorisFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			// NOW do the REVOKE on EVE
			userlib.DebugMsg("Alice revoking Eve's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentThree)))

			userlib.DebugMsg("Checking that Eve lost access to the file.")
			_, err = eve.LoadFile(eveFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that Eve cannot revoke the file she was revoked from")
			err = eve.RevokeAccess(eveFile, "alice")
			Expect(err).ToNot(BeNil())

			// Now check that frank adn grace can still append, load, etc
			err = frank.AppendToFile(frankFile, []byte(contentOne))
			Expect(err).To(BeNil())

			_, err = frank.LoadFile(frankFile)
			Expect(err).To(BeNil())

			_, err = grace.LoadFile(graceFile)
			Expect(err).To(BeNil())

			err = grace.AppendToFile(graceFile, []byte(contentOne))
			Expect(err).To(BeNil())

		})

		Specify("Testing the Constant Bandwidth Append", func() {

			// Reset bandwidth using DatastoreResetBandwidth()
			// 3) Using aliceLaptop, call AppendToFile(”test.txt”, ”abcdefghijklmnopqrstuvwxyz”)
			// 4) Define bandwidthAfter = DatastoreGetBandwidth().
			// 5) Check that bandwithAf ter < 50
			// 6) Define T = some really long string, Define S = sizeof(T) = 10000 bytes
			// 7) For file F in [”foo1.txt”, ”foo2.txt”, ”foo3.txt”, ..., ”foo100000.txt”] Alice calls StoreFile(F, T) on aliceLaptop
			// 8) Reset bandwidth using DatastoreResetBandwidth()
			// 9) Using aliceLaptop, AppendToFile(”foo50000.txt”, ”a”)
			// 10) Define bandwidthAfter = DatastoreGetBandwidth().
			// 11) Check that bandwidthAf ter < 50

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DatastoreResetBandwidth()

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			bandwidth1 := userlib.DatastoreGetBandwidth()

			userlib.DatastoreResetBandwidth()

			bandwidth2 := userlib.DatastoreGetBandwidth()

			bigContent := emptyString
			var i int
			for i = 1; i < 1000; i++ {
				bigContent = bigContent + contentTwo
			}

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(bigContent))
			Expect(err).To(BeNil())

			bandwidth3 := userlib.DatastoreGetBandwidth()

			userlib.DatastoreResetBandwidth()

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			bandwidth4 := userlib.DatastoreGetBandwidth()

			userlib.DebugMsg("%v %v %v %v", bandwidth1, bandwidth2, bandwidth3, bandwidth4)

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + bigContent + contentTwo)))
		})

		Specify("Testing invitation being tampered", func() {

			// initialize two users Alice and Bob and Mallory
			// Alice creates an invitation and send it to Bon
			// Mallary is trying to tamper it, using DataStoreGet access the invitation
			// Bob calls the acceptInvitation and should return an error about signature not valid

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initialing user Bob")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initialing mallicious user Mallory")
			// mallory, err := client.InitUser("mallory", defaultPassword)
			// Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DatastoreClear()

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("Testing Case Sensitive", func() {

			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			alice, err = client.InitUser("Alice", defaultPassword)
			Expect(err).To(BeNil())

		})

		Specify("Testing Empty password", func() {

			alice, err = client.InitUser("alice", "")
			Expect(err).To(BeNil())

			alice, err = client.InitUser("Alice", "")
			Expect(err).To(BeNil())

		})

		Specify("Testing Small Keystore", func() {

			alice, err = client.InitUser("alice", "")
			Expect(err).To(BeNil())

			userlib.KeystoreGetMap()

		})

		Specify("Testing Different Filenames", func() {

			alice, err = client.InitUser("alice", "sdfdsfsdf")
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile("abc.txt", []byte("sdfdsf"))
			Expect(err).To(BeNil())

			err = alice.StoreFile("sdf.txt", []byte("iji"))
			Expect(err).To(BeNil())

			_, err := alice.LoadFile("abc.txt")
			Expect(err).To(BeNil())

		})

		Specify("Testing WRONG password, + User Struct being tampered", func() {

			// initialize two users Alice and Bob and Mallory
			// Alice creates an invitation and send it to Bon
			// Mallary is trying to tamper it, using DataStoreGet access the invitation
			// Bob calls the acceptInvitation and should return an error about signature not valid

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			_, err = client.GetUser("alice", "garbage")
			Expect(err).ToNot(BeNil())

			userlib.DatastoreClear()

			_, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

		})

	})

})
