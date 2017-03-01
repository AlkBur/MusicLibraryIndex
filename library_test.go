package MusicLibraryIndex

import (
	"testing"
)

func TestBasicIndexBuilding(t *testing.T)  {
	library := New();
	library.AddTrack(&Track{
		Key:             "Anberlin/Never Take Friendship Personal/02. Paperthin Hymn.mp3",
		Name:            "Paperthin Hymn",
		ArtistName:      "Anberlin",
		AlbumName:       "Never Take Friendship Personal",
		Year:            2005,
		Genre:           "Other",
		Number:          2,
		AlbumArtistName: "Anberlin",
	})

	library.AddTrack(&Track{
		Key:             "Anberlin/Never Take Friendship Personal/08. The Feel Good Drag.mp3",
		Name:            "The Feel Good Drag",
		ArtistName:      "Anberlin",
		AlbumName:       "Never Take Friendship Personal",
		Year:            2005,
		Genre:           "Other",
		Number:          8,
		AlbumArtistName: "Anberlin",
	})

	library.RebuildTracks()

	track := library.GetTrack("Anberlin/Never Take Friendship Personal/08. The Feel Good Drag.mp3")
	strictEqual(track.Name, "The Feel Good Drag");
	strictEqual(track.ArtistName, "Anberlin");
	strictEqual(track.AlbumName, "Never Take Friendship Personal");
	strictEqual(track.Year, 2005);

	strictEqual(len(library.artistList), 1)
	artist := library.artistList[0]
	strictEqual(artist.Name, "Anberlin")
	strictEqual(artist.index, 0)
	strictEqual(len(artist.albumList), 1)
	album := artist.albumList[0];
	strictEqual(album.Name, "Never Take Friendship Personal");
	strictEqual(album.Year, 2005);
	strictEqual(album.index, 0);
	strictEqual(len(album.trackList), 2);
	strictEqual(album.trackList[0].Name, "Paperthin Hymn");
	strictEqual(album.trackList[0].index, 0);
	strictEqual(album.trackList[1].Name, "The Feel Good Drag");
	strictEqual(album.trackList[1].index, 1);



	strictEqual(len(library.albumList), 1);
	album = library.albumList[0];
	strictEqual(album.Name, "Never Take Friendship Personal");
	strictEqual(album.Year, 2005);
	strictEqual(album.index, 0);
	strictEqual(len(album.trackList), 2);
	strictEqual(album.trackList[0].Name, "Paperthin Hymn");
	strictEqual(album.trackList[0].index, 0);
	strictEqual(album.trackList[1].Name, "The Feel Good Drag");
	strictEqual(album.trackList[1].index, 1);



	var results = library.Search("never drag");

	strictEqual(len(results.albumList), 1);
	strictEqual(len(results.albumList[0].trackList), 1);
	strictEqual(results.albumList[0].trackList[0].Name, "The Feel Good Drag");

}

func TestCompilationAlbum(t *testing.T) {
	library := New()

	library.AddTrack(&Track{
		Key:             "jqvq-tpiu",
		Name:            "No News Is Good News",
		ArtistName:      "New Found Glory",
		AlbumName:       "2004 Warped Tour Compilation [Disc 1]",
		Year:            2004,
		disc:            1,
		discCount:       2,
		Genre:           "Alternative & Punk",
		AlbumArtistName: "Various Artists",
		track:           1,
	})

	library.AddTrack(&Track{
		Key:             "dldd-itve",
		Name:            "American Errorist (I Hate Hate Haters)",
		ArtistName:      "NOFX",
		AlbumName:       "2004 Warped Tour Compilation [Disc 1]",
		Year:            2004,
		disc:            1,
		discCount:       2,
		Genre:           "Alternative & Punk",
		AlbumArtistName: "Various Artists",
		track:           2,
	});
	library.AddTrack(&Track{
		Key:         "ukjv-ndsz",
		Name:        "Fire Down Below",
		ArtistName:  "Alkaline Trio",
		AlbumName:   "2007 Warped Tour Compilation [Disc 1]",
		compilation: true,
		Year:        2007,
		Genre:       "Alternative & Punk",
		track:       1,
		trackCount:  25,
	});
	library.AddTrack(&Track{
		Key:         "gfkt-esqz",
		Name:        "Requiem For Dissent",
		ArtistName:  "Bad Religion",
		AlbumName:   "2007 Warped Tour Compilation [Disc 1]",
		compilation: true,
		Year:        2007,
		Genre:       "Alternative & Punk",
		track:       2,
		trackCount:  25,
	});
	library.RebuildTracks()

	strictEqual(len(library.albumList), 2)
	strictEqual(len(library.artistList), 1)

	var artist = library.artistList[0]
	strictEqual(artist.Name, "Various Artists")
	strictEqual(len(artist.albumList), 2)
	strictEqual(len(library.albumList), 2)
	strictEqual(library.trackTable["jqvq-tpiu"].AlbumArtistName, "")
	strictEqual(library.trackTable["dldd-itve"].AlbumArtistName, "")
	strictEqual(library.trackTable["ukjv-ndsz"].AlbumArtistName, "")
	strictEqual(library.trackTable["gfkt-esqz"].AlbumArtistName, "")
}

func TestTracksFromSameAlbumMissingYearMetadata(t *testing.T) {
	library := New();
	library.AddTrack(&Track{
		Key:        "wwxj-unhr",
		Name:       "Dog-Eared Page",
		ArtistName: "The Matches",
		AlbumName:  "E. Von Dahl Killed the Locals",
		Year:       2004,
		Genre:      "Punk",
		track:      1,
	})

	library.AddTrack(&Track{
		Key:        "xekw-lvne",
		Name:       "Audio Blood",
		ArtistName: "The Matches",
		AlbumName:  "E. Von Dahl Killed the Locals",
		// missing year
		Genre: "Rock",
		track: 2,
	})

	library.AddTrack(&Track{
		Key:        "lpka-dugc",
		Name:       "Chain Me Free",
		ArtistName: "The Matches",
		AlbumName:  "E. Von Dahl Killed the Locals",
		Year:       2004,
		Genre:      "Rock",
		track:      3,
	});
	library.RebuildTracks()

	strictEqual(len(library.albumList), 1)
	strictEqual(library.albumList[0].Year, 2004)
	strictEqual(library.trackTable["xekw-lvne"].album.Year, 2004)
}

func TestDifferentAlbumsWithSameName(t *testing.T) {
	library := New()

	library.AddTrack(&Track{
		Key:        "sbao-lcvn",
		Name:       "6:00",
		ArtistName: "Dream Theater",
		AlbumName:  "Awake",
		Year:       1994,
		Genre:      "Progressive Rock",
		track:      1,
	})

	library.AddTrack(&Track{
		Key:        "qtru-gdtp",
		Name:       "Awake",
		ArtistName: "Godsmack",
		AlbumName:  "Awake",
		Year:       2000,
		Genre:      "Rock",
		track:      2,
	})

	library.RebuildTracks()

	strictEqual(len(library.albumList), 2)
}

func TestAlbumWithAFewTracksByDifferentArtists(t *testing.T) {
	library := New()

	library.AddTrack(&Track{
		Key:             "ikoe-nujf",
		Name:            "Paperthin Hymn",
		ArtistName:      "Anberlin",
		AlbumArtistName: "Anberlin",
		AlbumName:       "Never Take Friendship Personal",
		Year:            2005,
		Genre:           "Other",
		track:           2,
	})

	library.AddTrack(&Track{
		Key:             "msnq-swpc",
		Name:            "The Feel Good Drag",
		ArtistName:      "Anberlin, some other band",
		AlbumArtistName: "Anberlin",
		AlbumName:       "Never Take Friendship Personal",
		Year:            2005,
		Genre:           "Other",
		track:           8,
	})

	library.RebuildTracks()

	strictEqual(len(library.albumList), 1)
}

func TestAll(t *testing.T) {
	describe(t, "album by an artist", func() {
		library := New()

		library.AddTrack(&Track{
			Key:        "ynji-lcfu",
			Name:       "The Truth",
			ArtistName: "Relient K",
			AlbumName:  "Apathetic ep",
			track:      1,
			trackCount: 7,
		});
		library.AddTrack(&Track{
			Key:        "lxed-bsor",
			Name:       "Apathetic Way to Be",
			ArtistName: "Relient K",
			AlbumName:  "Apathetic ep",
			track:      2,
			trackCount: 7,
		});
		library.RebuildTracks();
		it(t, "should be filed under the artist", func() {
			strictEqual(len(library.artistList), 1);
			strictEqual(library.artistList[0].Name, "Relient K")
		})
	})

	describe(t, "album by an artist", func() {
		var library = New();
		library.AddTrack(&Track{
			Key:             "jqvq-tpiu",
			Name:            "No News Is Good News",
			ArtistName:      "New Found Glory",
			AlbumName:       "2004 Warped Tour Compilation",
			Year:            2004,
			disc:            1,
			discCount:       2,
			Genre:           "Alternative & Punk",
			AlbumArtistName: "Various Artists",
			track:           1,
		});
		library.AddTrack(&Track{
			Key:             "dldd-itve",
			Name:            "American Errorist (I Hate Hate Haters)",
			ArtistName:      "NOFX",
			AlbumName:       "2004 Warped Tour Compilation",
			Year:            2004,
			disc:            1,
			discCount:       2,
			Genre:           "Alternative & Punk",
			AlbumArtistName: "Various Artists",
			track:           2,
		})

		library.AddTrack(&Track{
			Key:         "ukjv-ndsz",
			Name:        "Fire Down Below",
			ArtistName:  "Alkaline Trio",
			AlbumName:   "2004 Warped Tour Compilation",
			disc:        2,
			compilation: true,
			Year:        2004,
			Genre:       "Alternative & Punk",
			track:       1,
			trackCount:  25,
		})

		library.AddTrack(&Track{
			Key:         "gfkt-esqz",
			Name:        "Requiem For Dissent",
			ArtistName:  "Bad Religion",
			AlbumName:   "2004 Warped Tour Compilation",
			disc:        2,
			compilation: true,
			Year:        2004,
			Genre:       "Alternative & Punk",
			track:       2,
			trackCount:  25,
		})

		library.RebuildTracks()
		it(t, "sorts by disc before track", func() {
			assertStrictEqual(t, library.albumList[0].trackList[0].Name, "No News Is Good News");
			assertStrictEqual(t, library.albumList[0].trackList[1].Name, "American Errorist (I Hate Hate Haters)");
			assertStrictEqual(t, library.albumList[0].trackList[2].Name, "Fire Down Below");
			assertStrictEqual(t, library.albumList[0].trackList[3].Name, "Requiem For Dissent");
		});
	})

	describe(t, "album artist with no album", func() {
		library := New();
		id := `5a89ea73-71aa-4c22-97a5-3b3509131cca`
		library.AddTrack(&Track{
			Key:             id,
			Name:            `I Miss You`,
			ArtistName:      `Blink 182`,
			ComposerName:    "",
			PerformerName:   "",
			AlbumArtistName: `blink-182`,
			AlbumName:       "",
			compilation:     false,
			track:           3,
			duration:        227.6815,
			Year:            2003,
			Genre:           `Rock`,
		})

		library.RebuildTracks();
		it(t, "shouldn't be various artists", func() {
			assertNotStrictEqual(t, library.trackTable[id].AlbumArtistName, "Various Artists");
		});
	})


	describe(t,"album with album artist", func() {
		var library = New();
		var id1 = `imnd-sxde`
		library.AddTrack(&Track{
			Key:             id1,
			Name:            `Palladio`,
			ArtistName:      `Escala`,
			AlbumArtistName: `Escala`,
			AlbumName:       `Escala`,
			track:           1,
		});
		var id2 = `vewu-hqbx`
		library.AddTrack(&Track{
			Key:             id2,
			Name:            `Requiem for a Tower`,
			ArtistName:      `Escala`,
			AlbumArtistName: `Escala`,
			AlbumName:       `Escala`,
			track:           2,
		})

		var id3 = `ixbc-oshh`
		library.AddTrack(&Track{
			Key:             id3,
			Name:            `Kashmir`,
			ArtistName:      `Escala; Slash`,
			AlbumArtistName: `Escala`,
			AlbumName:       `Escala`,
			track:           3,
		});
		library.RebuildTracks()

		it(t, "shouldn't be various artists", func() {
			assertStrictEqual(t, library.trackTable[id1].AlbumArtistName, "Escala");
			assertStrictEqual(t, library.trackTable[id2].AlbumArtistName, "Escala");
			assertStrictEqual(t, len(library.artistList), 1);
		})
	})

	describe(t, "label management", func() {
		var library = New()

		library.AddLabel(&Label{
			Id:   "wrong_id",
			Name: "wrong",
		});
		library.RebuildLabels();
		library.ClearLabels();
		library.AddLabel(&Label{
			Id:   "techno_id",
			Name: "techno",
		});
		library.AddLabel(&Label{
			Id:   "jazz_id",
			Name: "jazz",
		});
		library.RebuildLabels();
		it(t, "clearLabels, addLabel", func() {
			assertStrictEqual(t, library.labelList[0].Name, "jazz");
			assertStrictEqual(t, library.labelList[1].Name, "techno");
		});
	})

	describe(t,"searching with quoted seach terms", func() {
		var library = New()

		library.AddTrack(&Track{
			Key:        "fUPmxjMc",
			Name:       "Été (Original Mix)",
			ArtistName: "AKA AKA & Thalstroem",
			AlbumName:  "Varieté",
		})
		library.AddTrack(&Track{
			Key:        "zyGaKkrU",
			Name:       "Tribute to Young Stroke AKA Young Muscle",
			ArtistName: "Andy Kelley",
			AlbumName:  "The Weekend Challenge #3",
		})
		library.AddTrack(&Track{
			Key:        "v7zwEPLs",
			Name:       "Mista veri pakenee",
			ArtistName: "Turmion Katilot (no diacritics)",
			AlbumName:  "Pirun nyrkki",
		})
		library.AddTrack(&Track{
			Key:        "sobHcy0I",
			Name:       "Mistä veri pakenee",
			ArtistName: "Turmion Kätilöt (with diacritics)",
			AlbumName:  "Pirun nyrkki",
		})
		var literalQuoteKey = "G5FqXeJZ";
		library.AddTrack(&Track{
			Key:        literalQuoteKey,
			Name:       "A song with a literal \" in it",
			ArtistName: "Tester",
			AlbumName:  "literalQuote",
		})
		var literalBackslashKey = "Xsc4+ril";
		library.AddTrack(&Track{
			Key:        literalBackslashKey,
			Name:       "A song with a literal \\ in it",
			ArtistName: "Tester",
			AlbumName:  "literalBackslash",
		})

		library.RebuildTracks()

		it(t, "single search term returns both", func() {
			var results = library.Search("aka aka");
			assertStrictEqual(t, len(results.artistList), 2);
		})

		it(t, "quoted search term is case sensitive", func() {
			assertStrictEqual(t, len(library.Search("\"andy\"").artistList), 0);
			assertStrictEqual(t, len(library.Search("\"ANDY\"").artistList), 0);
			assertStrictEqual(t, len(library.Search("\"Andy\"").artistList), 1);
		});
		it(t, "quoted search terms include spaces", func() {
			var results = library.Search("\"AKA AKA\"");
			assertStrictEqual(t, len(results.artistList), 1)
			assertStrictEqual(t, results.artistList[0].Name, "AKA AKA & Thalstroem");
		});
		it(t, "quoted search terms preserve diacritics", func() {
			assertStrictEqual(t, len(library.Search("Mistä").artistList), 2);
			assertStrictEqual(t, len(library.Search("\"Mistä\"").artistList), 1);
		});
		it(t, "matches a song with a literal quote", func() {
			var results = library.Search("\"");
			assertStrictEqual(t, len(results.albumList), 1);
			assertStrictEqual(t, len(results.albumList[0].trackList), 1);
			assertStrictEqual(t, results.albumList[0].trackList[0].Key, literalQuoteKey);
		})
		it(t, "matches a song with a literal backslash", func() {
			var results = library.Search("\\");
			assertStrictEqual(t, len(results.albumList), 1);
			assertStrictEqual(t, len(results.albumList[0].trackList), 1);
			assertStrictEqual(t, results.albumList[0].trackList[0].Key, literalBackslashKey);
		})
	})

	describe(t, "searching with expressions", func() {
		var library = New();
		library.AddLabel(&Label{
			Id:   "techno_id",
			Name: "techno",
		})
		library.AddLabel(&Label{
			Id:   "jazz_id",
			Name: "jazz",
		});
		library.RebuildLabels()

		library.AddTrack(&Track{
			Key:        "fUPmxjMc",
			Name:       "Été (Original Mix)",
			ArtistName: "AKA AKA & Thalstroem",
			AlbumName:  "Varieté",
			Labels:     map[string]int{"jazz":1},
		});
		library.AddTrack(&Track{
			Key:        "v7zwEPLs",
			Name:       "Été (Remix)",
			ArtistName: "Some Remixer",
			AlbumName:  "Varieté",
			Labels:     map[string]int{"techno":1, "jazz":1},
		});
		library.AddTrack(&Track{
			Key:        "zyGaKkrU",
			Name:       "Tribute to Young Stroke AKA Young Muscle",
			ArtistName: "Andy Kelley",
			AlbumName:  "The Weekend Challenge #3",
		});
		library.RebuildTracks();

		it(t, "'not:'", func() {
			assertStrictEqual(t, len(library.Search(`not:andy`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`not:remix variete`).artistList), 1);
			assertStrictEqual(t, len(library.Search(`not:"AKA AKA"`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`not:"aka aka"`).artistList), 3);
			assertStrictEqual(t, len(library.Search(`not:(aka young)`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`not:`).artistList), 0);
		});
		it(t, "'label:'", func() {
			assertStrictEqual(t, len(library.Search(`label:techno`).artistList), 1);
			assertStrictEqual(t, len(library.Search(`label:jazz`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`not:label:techno`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`"label:techno"`).artistList), 0);
			assertStrictEqual(t, len(library.Search(`not:"label:techno"`).artistList), 3);
			assertStrictEqual(t, len(library.Search(`label:wrong`).artistList), 0);
			assertStrictEqual(t, len(library.Search(`not:label:wrong`).artistList), 3);
			assertStrictEqual(t, len(library.Search(`not:(label:jazz not:label:techno)`).artistList), 2);
		});
		it(t, "'or:'", func() {
			//assertStrictEqual(t, len(library.Search(`or:()`).artistList), 0);
			assertStrictEqual(t, len(library.Search(`or:(aka)`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`or:(variete)`).artistList), 2);
			assertStrictEqual(t, len(library.Search(`or:(label:techno not:label:jazz)`).artistList), 2);
		});
	})
}

func describe(t *testing.T, name string, f func())  {
	t.Log("TEST:", name)
	f()
}

func it(t *testing.T, name string, f func())  {
	t.Log("Equal:", name)
	f()
}

func assertStrictEqual(t *testing.T, src, dst interface{}) {
	if src != dst {
		t.Fatalf("Error: %v != %v", src, dst)
	}
}

func assertNotStrictEqual(t *testing.T, src, dst interface{}) {
	if src == dst {
		t.Fatalf("Error: %v != %v", src, dst)
	}
}