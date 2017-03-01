package MusicLibraryIndex

import (
	"strings"
	"sort"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
	"reflect"
)

type Library struct {
	searchFields       []string
	variousArtistsKey  string
	variousArtistsName string
	prefixesToStrip    []string

	trackTable  map[string]*Track
	artistTable map[string]*Artist
	artistList  ArtistList
	albumTable  map[string]*Album
	albumList   AlbumList
	dirtyTracks bool

	labelTable  map[string]*Label
	labelList   LabelList
	dirtyLabels bool
}

type AlbumList []*Album
type TrackList []*Track
type ArtistList []*Artist
type LabelList []*Label

type Track struct {
	Key             string
	Name            string
	ArtistName      string
	AlbumName       string
	AlbumArtistName string
	Year            int
	Genre           string
	Number          int

	ComposerName string
	PerformerName string

	Labels map[string]int

	//Album           *Album
	exactSearchTags string
	compilation     bool
	fuzzySearchTags string
	album           *Album
	//??????????
	disc  int
	discCount int
	track int
	index int
	trackCount int
	duration float64
}

type Album struct {
	Name  string
	Year  int
	Index int

	artist    *Artist
	trackList TrackList
	key       string
	index     int
}

type Artist struct {
	Name  string
	Index int

	albumList AlbumList
	key       string
	index     int
}

type Label struct {
	Id string
	Name string
	index int
}

type ArtistSet struct {
	val map[string][]bool
}

func NewArtistSet() *ArtistSet {
	return &ArtistSet{
		val: make(map[string][]bool),
	}
}

func (a *ArtistSet) Add(key string, val bool) {
	arr, ok := a.val[key]
	if !ok {
		arr := make([]bool, 1)
		arr[0] = val
		a.val[key] = arr
	} else {
		arr = append(arr, val)
	}
}

func (a *ArtistSet) MoreThanOneKey() bool {
	for _, val := range a.val {
		if len(val) > 1 {
			return true
		}
	}
	return false
}

func New() *Library {
	this := &Library{
		searchFields: []string{
			"ArtistName",
			"AlbumArtistName",
			"AlbumName",
			"Name",
		},
		variousArtistsKey:  "VariousArtists",
		variousArtistsName: "Various Artists",
		prefixesToStrip: []string{
			`/ ^\s*the\s+/`,
			`/^\s*a\s+/`,
			`/^\s*an\s+/`,
		},
	}
	this.ClearTracks()
	this.ClearLabels()
	return this
}

func (this *Library) ClearTracks() {
	this.trackTable = make(map[string]*Track)
	this.artistTable = make(map[string]*Artist)
	this.artistList = make([]*Artist, 0)
	this.albumTable = make(map[string]*Album)
	this.albumList = make([]*Album, 0)
	this.dirtyTracks = false;
}

func (this *Library) ClearLabels() {
	this.labelTable = make(map[string]*Label)
	this.labelList = make([]*Label, 0)
	this.dirtyLabels = false
}

func (this *Library) AddTrack(track *Track) {
	this.trackTable[track.Key] = track
	this.dirtyTracks = true
}

func (this *Library) GetTrack(key string) *Track {
	t, _ := this.trackTable[key]
	return t
}

func (this *Library) RebuildTracks() {
	if !this.dirtyTracks {
		return
	}
	this.RebuildAlbumTable()
	sort.Sort(this.albumList)

	var albumArtistName, artistKey string
	for albumKey := range this.albumTable {
		album := this.albumTable[albumKey]
		albumArtistSet := NewArtistSet()
		sort.Sort(album.trackList)
		albumArtistName = "";
		isCompilation := false;
		for i := 0; i < len(album.trackList); i += 1 {
			track := album.trackList[i];
			track.index = i;
			if track.AlbumArtistName == "" {
				albumArtistName = track.AlbumArtistName
				albumArtistSet.Add(this.getArtistKey(albumArtistName), true);
			}
			if albumArtistName == "" {
				albumArtistName = track.ArtistName
			}
			albumArtistSet.Add(this.getArtistKey(albumArtistName), true)
			isCompilation = track.compilation
		}
		if (isCompilation || albumArtistSet.MoreThanOneKey()) {
			albumArtistName = this.variousArtistsName
			artistKey = this.variousArtistsKey
			for i := 0; i < len(album.trackList); i += 1 {
				track := album.trackList[i]
				track.compilation = true
			}
		} else {
			artistKey = this.getArtistKey(albumArtistName)
		}
		artist := this.getOrCreateArtist(artistKey, albumArtistName)
		album.artist = artist;
		artist.albumList = append(artist.albumList, album)
	}
	this.artistList = make(ArtistList, 0)
	var variousArtist *Artist
	for artistKey := range this.artistTable {
		artist := this.artistTable[artistKey]
		//artist.albumList.sort(this.albumComparator);
		sort.Sort(artist.albumList)
		for i := 0; i < len(artist.albumList); i += 1 {
			album := artist.albumList[i];
			album.index = i;
		}
		if (artist.key == this.variousArtistsKey) {
			variousArtist = artist
		} else {
			this.artistList = append(this.artistList, artist)
		}
	}
	//this.artistList.sort(this.artistComparator);
	if variousArtist != nil {
		this.artistList = append(this.artistList, variousArtist)
	}
	sort.Sort(this.artistList)

	for i := 0; i < len(this.artistList); i += 1 {
		artist := this.artistList[i]
		artist.index = i;
	}

	this.dirtyTracks = false
}

func (this *Library) RebuildAlbumTable() {
	this.artistTable = make(map[string]*Artist)
	this.artistList = make([]*Artist, 0)
	this.albumTable = make(map[string]*Album)
	this.albumList = make([]*Album, 0)
	//thisAlbumList := this.albumList
	for trackKey := range this.trackTable {
		track := this.trackTable[trackKey]

		arrSearchTags := make([]string, 0, len(this.searchFields))
		refTrack := reflect.ValueOf(track)
		elTrack := refTrack.Elem()
		for _, field := range this.searchFields {
			f := elTrack.FieldByName(field)
			if f.IsValid() && f.CanSet() {
				str := f.String()
				if str != "" {
					arrSearchTags = append(arrSearchTags, str)
				}
			}
		}
		searchTags := strings.Join(arrSearchTags, "\n")

		track.exactSearchTags = searchTags
		track.fuzzySearchTags = formatSearchable(searchTags)
		if track.AlbumArtistName == this.variousArtistsName {
			track.AlbumArtistName = ""
			track.compilation = true
		}
		albumKey := this.getAlbumKey(track)
		album := this.getOrCreateAlbum(albumKey, track)
		track.album = album
		album.trackList = append(album.trackList, track)
		if (album.Year == 0) {
			album.Year = track.Year;
		}
	}
}

func (this *Library) getAlbumKey(track *Track) string {
	artistName := track.AlbumArtistName
	if artistName == "" {
		if track.compilation {
			artistName = this.variousArtistsName
		} else {
			artistName = track.ArtistName
		}
	}
	return formatSearchable(track.AlbumName + "\n" + artistName)
}

func (this *Library) getOrCreateAlbum(key string, track *Track) *Album {
	result, ok := this.albumTable[key]
	if !ok {
		result = &Album{
			Name:      track.AlbumName,
			Year:      track.Year,
			trackList: make(TrackList, 0),
			key:       key,
		}
		this.albumList = append(this.albumList, result)
		this.albumTable[key] = result
	}
	return result
}

func (this *Library) getOrCreateArtist(key, name string) *Artist {
	result, ok := this.artistTable[key]
	if !ok {
		result = &Artist{
			key:       key,
			Name:      name,
			albumList: make(AlbumList, 0),
		}
		this.artistTable[key] = result
	}
	return result
}

func (this *Library) getArtistKey(artistName string) string {
	return formatSearchable(artistName)
}

func (this *Library) Search(query string) *Library {
	searchResults := New()
	searchResults.searchFields = this.searchFields
	searchResults.variousArtistsKey = this.variousArtistsKey
	searchResults.variousArtistsName = this.variousArtistsName
	searchResults.prefixesToStrip = this.prefixesToStrip

	m := this.parseQuery(query);
	var track *Track
	for trackKey := range this.trackTable {
		track = this.trackTable[trackKey]
		if (m.matcher(track)) {
			searchResults.trackTable[track.Key] = track
		}
	}
	searchResults.dirtyTracks = true
	searchResults.RebuildTracks()

	return searchResults
}

func (this *Library) AddLabel(label *Label) {
	this.labelTable[label.Id] = label;
	this.dirtyLabels = true
}

func (this *Library) RebuildLabels() {
	if (!this.dirtyLabels) {
		return
	}

	this.labelList = make(LabelList, 0)
	for id := range this.labelTable {
		var label = this.labelTable[id]
		this.labelList = append(this.labelList, label)
	}

	sort.Sort(this.labelList)
	for index, label := range this.labelList {
		label.index = index;
	}

	this.dirtyLabels = false;
}

func formatSearchable(str string) string {
	return strings.ToLower(removeDiacritics(str))
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return result
}

func (this AlbumList) Len() int {
	return len(this)
}
func (this AlbumList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this AlbumList) Less(i, j int) bool {
	if this[i].Year < this[j].Year {
		return true
	} else if this[i].Year > this[j].Year {
		return false
	}
	return titleCompare(this[i].Name, this[j].Name)
}

func (this ArtistList) Len() int {
	return len(this)
}

func (this ArtistList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this ArtistList) Less(i, j int) bool {
	return titleCompare(this[i].Name, this[j].Name)
}

func (this TrackList) Len() int {
	return len(this)
}

func (this TrackList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this TrackList) Less(i, j int) bool {
	if this[i].disc < this[j].disc {
		return true
	} else if this[i].disc > this[j].disc {
		return false
	} else if this[i].track < this[j].track {
		return true
	} else if (this[i].track > this[j].track) {
		return false
	}
	return titleCompare(this[i].Name, this[j].Name)
}

func titleCompare(a, b string) bool {
	rez := strings.Compare(a, b)
	if rez < 0 {
		return true
	} else if rez > 0 {
		return false
	}
	return false
}

func (this LabelList) Len() int {
	return len(this)
}

func (this LabelList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this LabelList) Less(i, j int) bool {
	return titleCompare(this[i].Name, this[j].Name)
}
