package privkey

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	tcTestdataPath       = "./testdata/"
	tcTestdataDirPath    = "./testdata/directory"
	tcPrivKeyPath        = tcTestdataPath + "priv.key"
	tcInvalidPrivKeyPath = tcTestdataPath + "invalid_priv.key"
	tcTmpPrivKeyPath     = tcTestdataPath + "tmp_priv.key"

	tcPrivKey = `PrivKey{728B3595A3792E970780A6C3CFCB43B1BA108CBC3ECB9386E53302B7676BE2150EC3E1B601F780B2412FCDF11B0FE68469C4AF94BB5449BC18683658AD5A57F5DA4CB7A02B26627CC8A63869F5CC8BF29930087E9069902AA9B76AA21F3054246E9933EACA8759884A3FA6832C27176530B286EAB4D83383D0F83AC5E4303FBA2200252A6B11A47B489B560AB37D008F2DA00460799CFD2C11ECF60306444AC91C6CBE90A4CDF36F1AB4417E673535D79E2008337FAC3C0A026BECECC7FA86A73A53B15EC22237185907A30D4FB99F584298D8E39B21397108DA8091511F750C0B97CA6959487AEAB172C092AC82A6110A27280F7855BC62280F0B5382761F07262D678C5D1AFC6DC56C070960C5D01696994223C41CA7719A7F0F69C1283A9F62D8253BB9315F636169C83782E743CE60894E01834E9B3AA7538564200F0F350803C619527C98321B5796C3515E4CBE80F7018390A09066ACFF924D11F9A3E89A4553C01198A7BEEBB838AFF8467CF6C945422881E931D56730FEB171257790CAA7356177B8A4782128EA209E9A309F390F2BD517EBD4851C92BC2B8420F0731C1A3B26DFF04DDA87C5168178F0073B53DC63A670AFDCF16A4F18618B8A2E3D58409BEA8FC8A1769F93C320274D47F296C09C6F30AAA5905CB602913267D66CD9126D1EE7317664B4578B7022B76A23827046E5843EB0C815B88D16AB734A14B0220BA1922A2AEAE78DFAC738A5D48A9455B69797B25F42670B93BB89F60EF9B5CD755C35A4BCAF546952D243A30DB35EF7E91CF43C26703C57F8D673D5673E71F392CEF94665A2C7C9C3811824ABED3BBB0F135306353461A428E5E85E4AD4B9EFD31BBA138D579801F144AA2D74619EF46CBB511B4DDC8E52435C1F399830D909453C8891165AF0DB1EEE19B4BDBB71A2F67D17CAA3AD7463F422C4510A3C64A53A05D31D46EB668AF3AEC42B095157961A787C05D7267834C8EE1A270847BACDB39B95CC3C1837241F6102199B0D9D734B4ED00B39B24CD87380839CA872FCC9430B07B5E335DEF4742E288004AB5782AC21DEC8C7229B6BCD7839F69A20276C282D90964EF778B20810764A4DA629996E9607D77A36261867373A2707A298E5463B1C690FE7A05EEFE2AE82F664CCB08384791ACF426DBD56249A863DADE7B878121E6B3731B2172112B873F99B059BD34427042066369E24204129C4996DDAB9B3202B1DF4AB1B762546F8BF9DF8051FB8827C682C24BA23035205ABB49C3E174B04762BFDB6875E07178691B4DA169A58A45E6CEA8D73E38D46D8B3E7160CCF553B09EA8304C84EB6EC048586731CE838D8B64F0941BAE9B1687245C30006311CD470164B25B48476E0A141BB9C232E066EA3B1CF4850461A4746334B37C75CBD72B95F8B1346A369C39C01C8401135E1984F44B9B924C7976AE658B3B3C4FBFC1B03BA05FCB21F9FD11E9C750DB0AB47A0E9B571F595B6A5817A1758E4D5605CC03A6C36AA9CF720CBBB127562B3144B33C7E186766197417832C3D523D2718F0AE57EE4B1A1D729561BC801865C84BD09C652CC5693754C802243382B89B59C2C85F366AF0AD02540A63A217408538CB0A9A399449394E5CC4B01CD02EC397B1B2351A59941115FBB93268D84B807119166A091DFEA3D9CFA4028B1ADD5A98A0E651351539C2139C2A79B2FB54C32E6C874C5DA76EB2368DFDB4C7805BD7F0B3DAD9A944E6172A713144FB91837A71120F78B3016B599FA899EA606B839A7C2B81CCE23068A6A48428699005DB6DC2B570B619DACC2A05783C0E27B2906904220607AB10313934CAF3BB40B78590F482133E917828A904B11E14355E50B404B848524389041BD0B87C02079158F9C4928BB68E4D18D96C68E248C090E3A71AB71C5F10BAB6B25445179999B32723574B25AE8117F9B681DE76004155B5286A05183196838383CEA09AFA2B81D4757C4C3A04EBA7299628BDDE312B0F982C8E917CB5418E7E83D1E05144800C41CA544CF516D5AC2B448573B691240B52409DFBA57B1554C1658423DB67CB2F122D8BA2610293020A75B62497BED3B5F919ACF8B73C6FFA21F4A6B6DF5215E3364A32BAA5DC29212B951BBD77C8ED184AE4339C0626B3B22AB391671148464B0C3717B6DB1CF241A25CE018DD2CA428F634F1A16B091C573D7958649082DF6E0CB02880C4376B470263D2A5842FA4AA823043A62B9C9407B2604A147B543524CB427D7787477168B89AA71A092356DA11C37301176AA907E100BEAF4AAFEC1CC78860DC6BB3ECD1A6170E6A560935B0E1B5CE8A889EB41349D3152135AAFE08B8AD5DBB95692CFCA3A801884758EAB4C54A4BC18500FE167886056C07BE9B69C86009AB11AEDCC1066430E03640D9AB444683C93D61C6B3A862CF2351FB69CC2E1C183DD530892FB6668B67828112F86556C4AD1605CA0A579C62101CA8908C260B4D76349949E7038A21B90383432727694BD4A958242CB1E199B48274CBDB2FCACA729903D4A1160E6B87F0B9AA2A3BCCBB39353135E54E659ADCA16A8822F2E4214A128502C25573EC56D0AC8AEF123A5A159AFB5140C9279A45725A0CAC026FADB6CA8909234E86DADB228C4307B48258B7C885428CC909B88AEF495B2440B3C682870E6D795C6D10BCE11134EA94C4A5419F6056915D0A77A0C78FC09685A6430B5F894EDC050F744047AF382DF6B14B5823C354767EF2C4DAFF1B0F888B97AB48661D527A50703E553847C0BA148ABAF24F371F64BCFB1945F435056D3299C4D240708E1AC16065C11D557D804709BA86DF930445210B3FDB82DBED6A26774CDC442184D049128374388778A71188D79A1491EB32270016E4A243D47576073F64547E665A4166D44EA48FEAABAFBC006FD61065B24822879B536878AA8E9206ED230B8A4B67056BA42983C047282DBCC963108A0E7BB2CA055319965A5CA7B3CBC0176A4245009FC16171426E01ABED46C4FC56C456D7947BA48B0F1A808976954749580D9456A47B2818B362AF5C13A8C65C10DC364AE430740F40C324C8F01F9464B498CB4054EA9404EAD057DD61C8495A99155647D5162760AF3156948AEBA56451D94CC7517CC8EC9CBE107438653333377C267C808E9CC37469A8CF9A4138E4A1B1A6A3C5508A9C2F09CBC196F9A110B43290F4716A9BECB0F3C896648D7C53516C67B198D78B31A6B65BF6E824D5BF75049D998F19C41326A7EBB69B2D27934B999B6A2C1694F76BBBC596410CBBA4409DDA35CA1A8B289FE3FBA643150405F043D29A69F7A809D52BD57D4A4CCF15B8C8A4CD8A4D7DD5623FFEC021CC2CB7AE56774B63775152DB1880DBC612E0B3C1DF9232473D355BEC82B51BB7C2176C85711E3FC82C20CB0E842D22013DD33C35A5C3DE74F231DA8482AF824FF003A14E6CAB4FA73E3FB1FA30946A3CF712FC3F8DBE4C9D7108D7AFC84B8D3CE8DAA13D7DC7DB93F9BEDCBED8651F69DE5D662AACD0B35593B6727001668AA4B227DF91BAFD1FFE13793015A5359FDAFA0503FEBCEA7875D97E671AED3A46116EEC2B492C823C8D13B70EF170251423685047368123442470673360276417545103833882174302883010166550338356063020764360844744678266445872214663627428360158583030346165740074120467526840806566087485261223022866151612861676152106386272056301644384730684117635368258462712373152648270064450004081142463572617187815548880102445672524341383707174453654252476845781732631657701783647086188181584867618737172740842472842541013842508027815355701344852485382418005526581761361875527431760556725752472572367371048063480308577353425411873512387637765768732385437334462235434000165646522217552400506105386056162037033413164856148515770878627401438400502861313055018740557722628047748234200533183422451665406246582573240771607550001458585158021514376062400878201884636520463123741521582265657034701723020607870681858050764008160233380037341081432217415848712100658681814550534313484878602403825022864760680231801440335344146872258477230314124440876860514102773112706414264542882552548722211874806063038054763237567828257140335768588523313864808073610787258048656045301561002566666713608747202812040416187264738766375204626745261860770861263100778787426305252646737575744737308352852552265677408245377248845528134857432804156647811066370187775655611031608328777426540547748306253010885561731646222863543285626417835676122720481136278107627780370668083515850854854201858248463504132636765772660841030564883341275158103344670057352335760381042876034456083431327080685803871644665052712061108355636300853626640817866438086744881822687043642187744364863171558223432337632543550743558680766440347328563831631246827778323688354080057636077551701855114204167374323732763350380681100084103155650032400465060166381255485851256608748835255708145628441637653038056586701122746764071106454037376473113483821475738737007154223102366042482654685513726055312546570115386228484742161460671415565047504224654331328724408576861858305623512812534855736288346050653215310100188247734852504210471547550004830675826375462674461647733681023706438231751204663487505342612177326256025218258071270527455220638026284628730408477804040851146746072325447572065500422237228578746450018223323582620862877823267826187116517002748186078577443380110225213362810840718561150568707115271067634008047110800701001625468436187566066841110881154434413724716544252326184355788327088438250701157654166243820266671438187775130587071206123644624684403272535436523250860803410505065462187647126626170045648631816736705204843230617637278363471502861771617443133507328753256652777264330278305644410651616418202627247578121221382870554554571802357534742726572625230866671852660288521404625186161553607374748203646246513310660324152326371771241465171858382318116416857557184382612408482234780058215437400856608540043378741617116250380731755571155516331056423730565540268812373550731088873051553147721363632027075199DB0CB31543A7C6F152FF24E35DF9231D930B663A3B25251C916BAFABA38D5D00CFBCFBE88CED41B82F0CAA6D8DD618FD298E48ED4C1BD5BFE2C2D5988D8C6E12256D81FD71D5F9DAF98FBA7DF3253C43789F0DD6B8919E0A671F5F5AAAC8B7A458D8E6DC27C2CD635EAC509D3011EAED79193EB1952C5CEBBBB2D2DAA1C42501268F33F2B999A95FEB2C17B209D273F1FFF0D32F0D0C4E25181A46F4BAC5D1290C4701B5D52FDE9AA7B398C767991C22BC077A50D899AC8EA8C323843EFEAD9D8B7B1AF8ED633409B28B32326D88BFEF899365B0E5D09E587EA2B3EABA2A1387E79EEE8388397B2F3233A1E7599A8B6B99419164C6287FE24596B38DBA676D24BD7DFE04F9E35EC7032D3C491171E5FAAC902F224ADCC591ED934B2549BCB630DC29CD6F776B6D3D68660B3F4002AAEF2B751DA8586CF7C67527648EE41E1766AD295058D6CC20B80D23004B0FBD29AB892BC609240CB9F3F6EC70B1D0EB783D2F51E4ABF18CF82E7A6A19135D1AC14C665EE319B69554C9CF643678AB099EE4296C7BCF4CE6D1BC89D60307E43C618BB02965700A8FDB93FE7A1BD5BFB87F65841DD0A8C349465247CD94176F31055BBC8BDB1C6406E09E5DA93AE5303E5F6F8E916E030E8C976879C87F61BCE11090F7DEE82D91D14179B812B464532D89452712BE4C7A6598FAF81BDEADD5660A8378D42BAF059C73CB5E079733E8228738A120C3AD5EEE0F092FA1DA9FC1559BF5E04EB06BA071C5883479F4F00B939FD90C786D2B0D6F0E6E3656A339943658461AE2EF4644B68B8D8DEA0F538F562A31A563BDB5B2D3F5F2BAADC12C49D3856C813BB85C272BAA63AE7B400E641A3EF1AA3A6E8EE2AB7EB418D309DF1289584D6A4625663329017D41D5EFD6C99EA194B0683B2F14ECA3E38C6B23E00A42F7A28F7D9197F16412D1363424A9E6ED92AC09189C31B528771401329EF1684D6C0ACDC303AB25958BDAC5C7EA892B2448EA9DD3D9A47E881F319063A081BDA44893C2AB749583BB9891C4828D506B791708100E09DB6B6E132C4961210128F14F4367A0162CA8EEC94839FE622764B739D8CBE6839D3A6EEE6CCDD257112A15595CA01DAC75F3870FEA65ACE389B4A0634059615726B7641B30F11E50033D27087885CA6532AB20B71F4EC7B5AAC54C7E6FAEA88F3127C344298D554DA3E8F6DDC932E563AA754D971FEC633547CFE2ED7D761B3A68AB09C9FA296E97AC1CE15F8C3FCA380D8B2201A9231BCD481825B76B9D0A8F5D358030816AC7813720C8E44E4BD7013AECC5CE856BB7805D4D022181F9A0F800A7F9B4CFAD89A0D8CB61F95704E2C332F5F1A0944A96D3E6DE7C51142E2779824C9AF39C4A56AD7B9778A5EACBD6943B284D3A14ED35692E83186B54B90003890870702249472DC18CEE545806C6A43A76F02E30CCE22B25BC46C74037BFB5515AD6A0BD548DF4CF5314840FBBCD16B8EA91482213151A2215504D15B70D92D78D88DBEE6127D69E41F22CDB58E270EAAC1CA388112F3F343D755427F1EF473FEC3E1F61855FB4D936B8D45E73F9F32ABBD767039CF023536A90B1C447507A62EED9060CF307EB3CBB229F51EAED3B940B9FE4CE411E7DF9B5E3B39A36648679DA4C1AFA6F6FB57D15BF4536CFB63F652ABB5A0A85DCFDB1F1861C7D9DFAE5F3F155A5E4FF521A57D24DB85B9A61A70B8E987D5C86778B630F4C08177D45AA1B20138F5EDCE409E21206A41DBD889FEC97E6D294FD502F67F14CC97E40A5CE420BC5A782D557B530AB7FD9942BC96846A2E4908C0092EC54A802F8D981659663E7E486590810FFEB061B6417D8A9308757D437884D74E494570343FB68C191205767E3D0B9B28D091DD4198859CBE282535A764806C77FB6E9924E7FC64C7C8B32BE3002FD76312755A1874899CFD60FC3293715B0CA6C739CC017369265A09A709F01308FDFAEC35EBB61B76D56B70B64EA67347ABBF03720B3D518BDBD71863F32216DC1481B7ADD8C36564F7D6DB2FA0E11C2C9A6502B01F6E34D81858622577F2C26DF355D750617C7F3D39FDC13EF35EFE351948A06C2967170F1B7A68046A34C239C11CF6ECA1C2840E19B9F55B56F3AD2EE81FC4E3B6307A74FE4711E3F099956782F1338D4008C8376C400D233EDF0CF3FA02150FFC4C56678E4D4C4419EBC85F0E0CF5365FB50834961B76A5F46DB8D576A5F70E77F0DEE164D3AD4A461BE2E14A89288C4D2DA180530CEA476CD90A4F3EFCAE1C6E1FF98D33868BF186FFF121F55D1DBBD74E89AC2F935DFDB59CBF4C6C8AC49C47C31C64FB095751442AA532DB16234BEC2294B3495D7298EF363CA3B5333DD3BF57B18174EB11BCD2B23AD22DA5F59DF7992AFB913F0D8E1E2CC5ED9FE46E1990C8A9671510BF05BEF793F4B8AD9A59AB8CA83736A04E4178427DFBA0423DEC209ADC62A8D3D767B3562B7849A45A5FC026D6AE87230948F8D5B957F4A05A785DE93E1C8A6748EFC0B20E80E8ED56C8B10BE260E68E730E3264F6624272BDFD29380CFCDA7E98B8B1BA6FD9DCC2728DCFF7D901B874071D5393F41E835D92DDBE931426DBE85831EC26098185E1729A3277E00577C81A61C06448F5771DC101FCC93B5E54BBE19A6B6AC0EF59B826D2ADCAFA2080783CC7DB7E3A979340C019FE3D2688D85DA3F5DC139A026B3A67B90DEB46F591A03C063CF2D81914A1BF20588E71FB8E45F465B00A1C65DAABD271C6BACFD644F8F75EA9923694882C61B88BCBFD3C00152BDD20CED3904EBE9B4FC85194580DB548F7500B0B5DB750FA0985C643108EDAD3C4B10695B10B853BECD11B395D16A3FE0C190DC9FE74383C8F3DC6976E58B286C85932D96734316009AA66364484DAEEC1F7E514693E89B6E1B859E6EB77276990B7748AF88E131BC34F9610548B4EDF9728A2E60551D4FAE81232F57224EE4EAB37C37F781063B2577412D65CF34304B4A4D7F11BAD8CF9C51DC78B155CA3594369C3880DFD1182DB4D7095DAFFC7B5F2F0B062943E0DFC1E6343C3588727B198EAE3C295CAA143797D2B091BE6851F6536BE24D9C804064089EFA24FA3AA042DA1532393C2B0F4AEE265C2235202E8ED2DC048BC2886DF804C039115F3227E22A9E1FAEF1511C3B5CF6C49784AF98F2A4D7BAF87031A9AA52479250637D942424CECC81216403B17FC53C772E50E072D4576EDA1906389EB44D8D46350B7BB8D24046861A0EFC9E812F2A707D375DC2C5C761674EFE06FD01037A67FEC590BD7FFA23A37C25BB83C7DDAF23BAEC0095BA1B64AB7B29734004D2B7789D105559A6C79309DA9DB5B95DD0FE9B1F2E862A93E6E6220D0553E59FCEF18D585B692212C194037B8B16323B89F3C95A6D4B76B122F2E9E7E8D5E99741951D7272993FCAE8646870637CC635F3C6BB3E075D7C2794162481541B5A770776D757A004179DBC2EC0C591FF455FFC81A59E7DE4E8B2482571E92CEA741A62866DC11997E7E6189305B8303248CA674129047DA5CB78E914DF32249546059DBD}`
)

func testDeleteFile(f string) {
	os.RemoveAll(f)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestGetPrivKey(t *testing.T) {
	t.Parallel()

	testDeleteFile(tcTmpPrivKeyPath)
	defer testDeleteFile(tcTmpPrivKeyPath)

	privKey, err := GetPrivKey(tcPrivKeyPath)
	if err != nil {
		t.Error(err)
		return
	}
	if privKey.ToString() != asymmetric.LoadPrivKey(tcPrivKey).ToString() {
		t.Error("diff private keys")
		return
	}

	if _, err := GetPrivKey(tcInvalidPrivKeyPath); err == nil {
		t.Error("success get invalid private key")
		return
	}
	if _, err := GetPrivKey("./random/not_exist/path/57199u140291724y121291d1/priv.key"); err == nil {
		t.Error("success get private key with not exist directory")
		return
	}
	if _, err := GetPrivKey(tcTestdataDirPath); err == nil {
		t.Error("success get private key as directory")
		return
	}

	tmpPrivKey, err := GetPrivKey(tcTmpPrivKeyPath)
	if err != nil {
		t.Error(err)
		return
	}
	tmpPrivKeyX, err := GetPrivKey(tcTmpPrivKeyPath)
	if err != nil {
		t.Error(err)
		return
	}
	if tmpPrivKey.ToString() != tmpPrivKeyX.ToString() {
		t.Error("diff tmp private keys")
		return
	}
}
