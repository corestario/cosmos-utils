package storewrapper_test

import (
	"io/ioutil"
	"storewrapper"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func makeTestStore() (*storewrapper.KVStore, error) {
	testKey := sdk.NewKVStoreKey("testKey")

	dbDir, err := ioutil.TempDir("", "goleveldb-app-sim")
	if err != nil {
		return nil, err
	}

	db, err := sdk.NewLevelDB("Simulation", dbDir)
	if err != nil {
		return nil, err
	}

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testKey, sdk.StoreTypeIAVL, db)
	err = ms.LoadLatestVersion()
	if err != nil {
		return nil, err
	}
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	return storewrapper.NewKVStore(ctx.KVStore(testKey), 0), nil

}

func TestKVStore(t *testing.T) {
	stw, err := makeTestStore()
	require.Nil(t, err)

	key := []byte("Ezekiel")
	parts, err := stw.SetW(key, []byte(sampleLarge))
	require.Nil(t, err)
	require.NotZero(t, parts)

	n1, _, err := stw.HasW(key)
	require.Nil(t, err)
	require.NotZero(t, n1)
	require.Equal(t, n1, parts)

	res, err := stw.GetW(key)
	require.Nil(t, err)
	require.Equal(t, len(sampleLarge), len(res))
	parts, err = stw.SetW(key, []byte(sampleSmall))
	require.Nil(t, err)
	require.NotZero(t, parts)

	n2, _, err := stw.HasW(key)
	require.Nil(t, err)
	require.NotZero(t, n2)
	require.NotEqual(t, n2, n1)
	require.Equal(t, n2, parts)
}

var sampleSmall = `17 And I will execute great vengeance upon them with furious rebukes; and they shall know that I am the Lord, when I shall lay my vengeance upon them.`

var sampleLarge = `Ezekiel 25 King James Version (KJV)
25 The word of the Lord came again unto me, saying,

2 Son of man, set thy face against the Ammonites, and prophesy against them;

3 And say unto the Ammonites, Hear the word of the Lord God; Thus saith the Lord God; Because thou saidst, Aha, against my sanctuary, when it was profaned; and against the land of Israel, when it was desolate; and against the house of Judah, when they went into captivity;

4 Behold, therefore I will deliver thee to the men of the east for a possession, and they shall set their palaces in thee, and make their dwellings in thee: they shall eat thy fruit, and they shall drink thy milk.

5 And I will make Rabbah a stable for camels, and the Ammonites a couching place for flocks: and ye shall know that I am the Lord.

6 For thus saith the Lord God; Because thou hast clapped thine hands, and stamped with the feet, and rejoiced in heart with all thy despite against the land of Israel;

7 Behold, therefore I will stretch out mine hand upon thee, and will deliver thee for a spoil to the heathen; and I will cut thee off from the people, and I will cause thee to perish out of the countries: I will destroy thee; and thou shalt know that I am the Lord.

8 Thus saith the Lord God; Because that Moab and Seir do say, Behold, the house of Judah is like unto all the heathen;

9 Therefore, behold, I will open the side of Moab from the cities, from his cities which are on his frontiers, the glory of the country, Bethjeshimoth, Baalmeon, and Kiriathaim,

10 Unto the men of the east with the Ammonites, and will give them in possession, that the Ammonites may not be remembered among the nations.

11 And I will execute judgments upon Moab; and they shall know that I am the Lord.

12 Thus saith the Lord God; Because that Edom hath dealt against the house of Judah by taking vengeance, and hath greatly offended, and revenged himself upon them;

13 Therefore thus saith the Lord God; I will also stretch out mine hand upon Edom, and will cut off man and beast from it; and I will make it desolate from Teman; and they of Dedan shall fall by the sword.

14 And I will lay my vengeance upon Edom by the hand of my people Israel: and they shall do in Edom according to mine anger and according to my fury; and they shall know my vengeance, saith the Lord God.

15 Thus saith the Lord God; Because the Philistines have dealt by revenge, and have taken vengeance with a despiteful heart, to destroy it for the old hatred;

16 Therefore thus saith the Lord God; Behold, I will stretch out mine hand upon the Philistines, and I will cut off the Cherethims, and destroy the remnant of the sea coast.

17 And I will execute great vengeance upon them with furious rebukes; and they shall know that I am the Lord, when I shall lay my vengeance upon them.


Ezekiel 26 King James Version (KJV)

26 And it came to pass in the eleventh year, in the first day of the month, that the word of the Lord came unto me, saying,

2 Son of man, because that Tyrus hath said against Jerusalem, Aha, she is broken that was the gates of the people: she is turned unto me: I shall be replenished, now she is laid waste:

3 Therefore thus saith the Lord God; Behold, I am against thee, O Tyrus, and will cause many nations to come up against thee, as the sea causeth his waves to come up.

4 And they shall destroy the walls of Tyrus, and break down her towers: I will also scrape her dust from her, and make her like the top of a rock.

5 It shall be a place for the spreading of nets in the midst of the sea: for I have spoken it, saith the Lord God: and it shall become a spoil to the nations.

6 And her daughters which are in the field shall be slain by the sword; and they shall know that I am the Lord.

7 For thus saith the Lord God; Behold, I will bring upon Tyrus Nebuchadrezzar king of Babylon, a king of kings, from the north, with horses, and with chariots, and with horsemen, and companies, and much people.

8 He shall slay with the sword thy daughters in the field: and he shall make a fort against thee, and cast a mount against thee, and lift up the buckler against thee.

9 And he shall set engines of war against thy walls, and with his axes he shall break down thy towers.

10 By reason of the abundance of his horses their dust shall cover thee: thy walls shall shake at the noise of the horsemen, and of the wheels, and of the chariots, when he shall enter into thy gates, as men enter into a city wherein is made a breach.

11 With the hoofs of his horses shall he tread down all thy streets: he shall slay thy people by the sword, and thy strong garrisons shall go down to the ground.

12 And they shall make a spoil of thy riches, and make a prey of thy merchandise: and they shall break down thy walls, and destroy thy pleasant houses: and they shall lay thy stones and thy timber and thy dust in the midst of the water.

13 And I will cause the noise of thy songs to cease; and the sound of thy harps shall be no more heard.

14 And I will make thee like the top of a rock: thou shalt be a place to spread nets upon; thou shalt be built no more: for I the Lord have spoken it, saith the Lord God.

15 Thus saith the Lord God to Tyrus; Shall not the isles shake at the sound of thy fall, when the wounded cry, when the slaughter is made in the midst of thee?

16 Then all the princes of the sea shall come down from their thrones, and lay away their robes, and put off their broidered garments: they shall clothe themselves with trembling; they shall sit upon the ground, and shall tremble at every moment, and be astonished at thee.

17 And they shall take up a lamentation for thee, and say to thee, How art thou destroyed, that wast inhabited of seafaring men, the renowned city, which wast strong in the sea, she and her inhabitants, which cause their terror to be on all that haunt it!

18 Now shall the isles tremble in the day of thy fall; yea, the isles that are in the sea shall be troubled at thy departure.

19 For thus saith the Lord God; When I shall make thee a desolate city, like the cities that are not inhabited; when I shall bring up the deep upon thee, and great waters shall cover thee;

20 When I shall bring thee down with them that descend into the pit, with the people of old time, and shall set thee in the low parts of the earth, in places desolate of old, with them that go down to the pit, that thou be not inhabited; and I shall set glory in the land of the living;

21 I will make thee a terror, and thou shalt be no more: though thou be sought for, yet shalt thou never be found again, saith the Lord God.
`
