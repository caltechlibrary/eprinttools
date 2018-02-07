//
// Package eprinttools is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package eprinttools

import (
	"encoding/xml"
	"net/url"
	"os"
	"testing"

	// Caltech Library Packages
	"github.com/caltechlibrary/rc"
)

func TestEPrint3x(t *testing.T) {
	// Simulate URL response for https://authors.library.caltech.edu/rest/eprint/84590.xml
	src := []byte(`<?xml version='1.0' encoding='utf-8'?>
<eprints xmlns='http://eprints.org/ep2/data/2.0'>
  <eprint id='https://authors.library.caltech.edu/id/eprint/84590'>
    <eprintid>84590</eprintid>
    <rev_number>10</rev_number>
    <documents>
      <document id='https://authors.library.caltech.edu/id/document/263914'>
        <docid>263914</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/1374695'>
            <fileid>1374695</fileid>
            <datasetid>document</datasetid>
            <objectid>263914</objectid>
            <filename>PhysRevB.97.014311.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <hash>1caf0b49dba8b63fd36e56218ade5075</hash>
            <hash_type>MD5</hash_type>
            <filesize>1900134</filesize>
            <mtime>2018-01-30 22:32:46</mtime>
            <url>https://authors.library.caltech.edu/84590/1/PhysRevB.97.014311.pdf</url>
          </file>
        </files>
        <eprintid>84590</eprintid>
        <pos>1</pos>
        <placement>1</placement>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <language>en</language>
        <security>public</security>
        <license>other</license>
        <main>PhysRevB.97.014311.pdf</main>
        <content>published</content>
      </document>
      <document id='https://authors.library.caltech.edu/id/document/263915'>
        <docid>263915</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/1374697'>
            <fileid>1374697</fileid>
            <datasetid>document</datasetid>
            <objectid>263915</objectid>
            <filename>1710.09843.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <hash>2cc920012716ac8f74bee91a882dfc1e</hash>
            <hash_type>MD5</hash_type>
            <filesize>2371729</filesize>
            <mtime>2018-01-30 22:32:52</mtime>
            <url>https://authors.library.caltech.edu/84590/2/1710.09843.pdf</url>
          </file>
        </files>
        <eprintid>84590</eprintid>
        <pos>2</pos>
        <placement>2</placement>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <language>en</language>
        <security>public</security>
        <license>other</license>
        <main>1710.09843.pdf</main>
        <content>submitted</content>
      </document>
    </documents>
    <eprint_status>archive</eprint_status>
    <userid>18</userid>
    <dir>disk0/00/08/45/90</dir>
    <datestamp>2018-01-31 00:00:27</datestamp>
    <lastmod>2018-01-31 00:00:27</lastmod>
    <status_changed>2018-01-31 00:00:27</status_changed>
    <type>article</type>
    <metadata_visibility>show</metadata_visibility>
    <creators>
      <item>
        <name>
          <family>Seetharam</family>
          <given>Karthik</given>
        </name>
        <id>Seetharam-K-I</id>
      </item>
      <item>
        <name>
          <family>Titum</family>
          <given>Paraj</given>
        </name>
        <id>Titum-P</id>
      </item>
      <item>
        <name>
          <family>Kolodrubetz</family>
          <given>Michael</given>
        </name>
        <id>Kolodrubetz-M</id>
      </item>
      <item>
        <name>
          <family>Refael</family>
          <given>Gil</given>
        </name>
        <id>Refael-G</id>
      </item>
    </creators>
    <title>Absence of thermalization in finite isolated interacting Floquet systems</title>
    <ispublished>pub</ispublished>
    <full_text_status>public</full_text_status>
    <note>© 2018 American Physical Society. 

Received 8 November 2017; revised manuscript received 6 January 2018; published 29 January 2018. 

The authors would like to thank J. Garrison, M. Bukov,
A. Polkovnikov, E. van Nieuwenburg, Y. Baum, M.-F. Tu,
and J. Moore for insightful discussions. M.K. acknowledges
funding from Laboratory Directed Research and Development
from Berkeley Laboratory, provided by the Director, Office
of Science, of the US Department of Energy under Contract
No. DEAC02-05CH11231, and from the US DOE, Office
of Science, Basic Energy Sciences, as part of the TIMES
initiative. G.R. and K.S. are grateful for support from the NSF through DMR-1410435, the Institute of Quantum Information and Matter, an NSF Frontier center funded by the Gordon and Betty Moore Foundation, the Packard Foundation, and from the ARO MURI W911NF-16-1-0361 “Quantum Materials by Design with Electromagnetic Excitation” sponsored by the US Army. K.S. is additionally grateful for support from NSF Graduate Research Fellowship Program. P.T. is supported by a National Research Council postdoctoral fellowship, and acknowledges funding from ARL CDQI, NSF PFC at JQI, ARO, AFOSR, ARO MURI, and NSF QIS.</note>
    <abstract>Conventional wisdom suggests that the long-time behavior of isolated interacting periodically driven (Floquet) systems is a featureless maximal-entropy state characterized by an infinite temperature. Efforts to thwart this uninteresting fixed point include adding sufficient disorder to realize a Floquet many-body localized phase or working in a narrow region of drive frequencies to achieve glassy nonthermal behavior at long time. Here we show that in clean systems the Floquet eigenstates can exhibit nonthermal behavior due to finite system size. We consider a one-dimensional system of spinless fermions with nearest-neighbor interactions where the interaction term is driven. Interestingly, even with no static component of the interaction, the quasienergy spectrum contains gaps and a significant fraction of the Floquet eigenstates, at all quasienergies, have nonthermal average doublon densities. We show that this nonthermal behavior arises due to emergent integrability at large interaction strength and discuss how the integrability breaks down with power-law dependence on system size.</abstract>
    <date>2018-01-01</date>
    <date_type>published</date_type>
    <publication>Physical Review B</publication>
    <volume>97</volume>
    <number>1</number>
    <publisher>American Physical Society</publisher>
    <pagerange>Art. No. 014311</pagerange>
    <id_number>CaltechAUTHORS:20180130-143240205</id_number>
    <refereed>TRUE</refereed>
    <issn>2469-9950</issn>
    <official_url>http://resolver.caltech.edu/CaltechAUTHORS:20180130-143240205</official_url>
    <related_url>
      <item>
        <url>https://doi.org/10.1103/PhysRevB.97.014311</url>
        <type>doi</type>
        <description>Article</description>
      </item>
      <item>
        <url>https://journals.aps.org/prb/abstract/10.1103/PhysRevB.97.014311</url>
        <type>pub</type>
        <description>Article</description>
      </item>
      <item>
        <url>https://arxiv.org/abs/1710.09843</url>
        <type>arxiv</type>
        <description>Discussion Paper</description>
      </item>
    </related_url>
    <referencetext>
      <item>T. Oka and H. Aoki, Phys. Rev. B 79, 081406 (2009).
T. Kitagawa, E. Berg, M. Rudner, and E. Demler, Phys. Rev. B 82, 235114 (2010).
N. H. Lindner, G. Refael, and V. Galitski, Nat. Phys. 7, 490 (2011).
L. Jiang, T. Kitagawa, J. Alicea, A. R. Akhmerov, D. Pekker, G. Refael, J. I. Cirac, E. Demler, M. D. Lukin, and P. Zoller, Phys. Rev. Lett. 106, 220402 (2011).
T. Kitagawa, T. Oka, A. Brataas, L. Fu, and E. Demler, Phys. Rev. B 84, 235108 (2011).
V. Khemani, A. Lazarides, R. Moessner, and S. L. Sondhi, Phys. Rev. Lett. 116, 250401 (2016).
D. V. Else, B. Bauer, and C. Nayak, Phys. Rev. Lett. 117, 090402 (2016).
N. Y. Yao, A. C. Potter, I.-D. Potirniche, and A. Vishwanath, Phys. Rev. Lett. 118, 030401 (2017).
C. W. von Keyserlingk and S. L. Sondhi, Phys. Rev. B 93, 245145 (2016).
C. W. von Keyserlingk and S. L. Sondhi, Phys. Rev. B 93, 245146 (2016).
M. S. Rudner, N. H. Lindner, E. Berg, and M. Levin, Phys. Rev. X 3, 031005 (2013).
P. Titum, E. Berg, M. S. Rudner, G. Refael, and N. H. Lindner, Phys. Rev. X 6, 021013 (2016).
H. C. Po, T. Fidkowski, A. C. Potter, and A. Vishwanath, Phys. Rev. B 96, 245116 (2017).
D. V. Else and C. Nayak, Phys. Rev. B 93, 201103 (2016).
A. C. Potter, T. Morimoto, and A. Vishwanath, Phys. Rev. X 6, 041001 (2016).
R. Roy and F. Harper, Phys. Rev. B 94, 125105 (2016).
H. C. Po, L. Fidkowski, T. Morimoto, A. C. Potter, and A. Vishwanath, Phys. Rev. X 6, 041070 (2016).
F. Harper and R. Roy, Phys. Rev. Lett. 118, 115301 (2017).
A. C. Potter and T. Morimoto, Phys. Rev. B 95, 155126 (2017).
R. Roy and F. Harper, Phys. Rev. B 95, 195128 (2017).
A. C. Potter, A. Vishwanath, and L. Fidkowski, , , and , arXiv:1706.01888.
J. Zhang, P. W. Hess, A. Kyprianidis, P. Becker, A. Lee, J. Smith, G. Pagano, I.-D. Potirniche, A. C. Potter, and A. Vishwanath et al., Nature (London) 543, 217 (2017).
S. Choi, J. Choi, R. Landig, G. Kucsko, H. Zhou, J. Isoya, F. Jelezko, S. Onoda, H. Sumiya, and V. Khemani et al., Nature (London) 543, 221 (2017).
G. Jotzu, M. Messer, R. Desbuquois, M. Lebrat, T. Uehlinger, D. Greif, and T. Esslinger, Nature (London) 515, 237 (2014).
M. C. Rechtsman, J. M. Zeuner, Y. Plotnik, Y. Lumer, D. Podolsky, F. Dreisow, S. Nolte, M. Segev, and A. Szameit, Nature (London) 496, 196 (2013).
P. Bordia, H. Luschen, U. Schneider, M. Knap, and I. Bloch, Nat. Phys. 13, 460 (2017).
M. Aidelsburger, M. Atala, M. Lohse, J. T. Barreiro, B. Paredes, and I. Bloch, Phys. Rev. Lett. 111, 185301 (2013).
H. Miyake, G. A. Siviloglou, C. J. Kennedy, W. C. Burton, and W. Ketterle, Phys. Rev. Lett. 111, 185302 (2013).
C. Parker, L.-C. Ha, and C. Chin, Nat. Phys. 9, 769 (2013).
R. Nandkishore and D. A. Huse, Annu. Rev. Condens. Matter Phys. 6, 15 (2015).
L. D&apos;Alessio, Y. Kafri, A. Polkovnikov, and M. Rigol, Adv. Phys. 65, 239 (2016).
J. M. Deutsch, Phys. Rev. A 43, 2046 (1991).
M. Srednicki, Phys. Rev. E 50, 888 (1994).
M. Rigol, V. Dunjko, and M. Olshanii, Nature (London) 452, 854 (2008).
L. D&apos;Alessio and M. Rigol, Phys. Rev. X 4, 041048 (2014).
A. Lazarides, A. Das, and R. Moessner, Phys. Rev. Lett. 112, 150401 (2014).
A. Lazarides, A. Das, and R. Moessner, Phys. Rev. E 90, 012110 (2014).
K. I. Seetharam, C.-E. Bardyn, N. H. Lindner, M. S. Rudner, and G. Refael, Phys. Rev. X 5, 041050 (2015).
T. Iadecola, T. Neupert, and C. Chamon, Phys. Rev. B 91, 235133 (2015).
T. Iadecola and C. Chamon, Phys. Rev. B 91, 184301 (2015).
H. Dehghani, T. Oka, and A. Mitra, Phys. Rev. B 90, 195429 (2014).
L. Vidmar and M. Rigol, J. Stat. Mech.: Theory Exp. 2016, 064007.
A. Polkovnikov, K. Sengupta, A. Silva, and M. Vengalattore, Rev. Mod. Phys. 83, 863 (2011).
D. A. Abanin and Z. Papić, Ann. Phys. 529, 1700169 (2017).
J. R. Garrison, R. V. Mishmash, and M. P. A. Fisher, Phys. Rev. B 95, 054204 (2017).
N. Y. Yao, C. R. Laumann, J. I. Cirac, M. D. Lukin, and J. E. Moore, Phys. Rev. Lett. 117, 240601 (2016).
A. Smith, J. Knolle, D. L. Kovrizhin, and R. Moessner, Phys. Rev. Lett. 118, 266601 (2017).
L. D&apos;Alessio and A. Polkovnikov, Ann. Phys. 333, 19 (2013).
P. Ponte, A. Chandran, Z. Papić, and D. A. Abanin, Ann. Phys. 353, 196 (2014).
P. Ponte, Z. Papić, F. Huveneers, and D. A. Abanin, Phys. Rev. Lett. 114, 140401 (2015).
A. Lazarides, A. Das, and R. Moessner, Phys. Rev. Lett. 115, 030402 (2015).
A. Agarwala and D. Sen, Phys. Rev. B 95, 014305 (2017).
S. A. Weidinger and M. Knap, Sci. Rep. 7, 45382 (2017).
F. Machado, G. D. Meyer, D. V. Else, C. Nayak, and N. Y. Yao, , , , , and , arXiv:1708.01620.
M. Bukov, S. Gopalakrishnan, M. Knap, and E. Demler, Phys. Rev. Lett. 115, 205301 (2015).
D. V. Else, B. Bauer, and C. Nayak, Phys. Rev. X 7, 011026 (2017).
T.-S. Zeng and D. N. Sheng, Phys. Rev. B 96, 094202 (2017).
D. Abanin, W. De Roeck, W. W. Ho, and F. Huveneers, Commun. Math. Phys. 354, 809 (2017).
T. Kuwahara, T. Mori, and K. Saito, Ann. Phys. 367, 96 (2016).
J. Sirker and M. Bortz, J. Stat. Mech.: Theory Exp. 2006, P01007.
M. Rigol and L. F. Santos, Phys. Rev. A 82, 011604 (2010).
L. F. Santos and M. Rigol, Phys. Rev. E 82, 031130 (2010).
P. W. Claeys and J.-S. Caux, and , arXiv:1708.07324.
Although, of course, the precise scaling and crossover behavior indeed should depend on the model.
M. Bukov, M. Heyl, D. A. Huse, and A. Polkovnikov, Phys. Rev. B 93, 155132 (2016).
M. Bukov, M. Kolodrubetz, and A. Polkovnikov, Phys. Rev. Lett. 116, 125301 (2016).
J. H. Mentink, K. Balzer, and M. Eckstein, Nat. Commun. 6, 6708 (2015).
F. Görg, M. Messer, K. Sandholzer, G. Jotzu, R. Desbuquois, and T. Esslinger, , , , , , and , arXiv:1708.06751.
S. Kitamura, T. Oka, and H. Aoki, Phys. Rev. B 96, 014406 (2017).
W. Nie, H. Katsura, and M. Oshikawa, Phys. Rev. Lett. 111, 100402 (2013).
We assume open boundary conditions throughout the enire paper. In this case, there is no distinction between fermions and hard-core bosons. However, if considering the system on a ring, then one must be careful about (anti)periodic boundary conditions as exchange statistics are relevant [70].
This is achievable using the probability of finding a 
k
 -doublon state of a half-filled 
N
 site system given by 
p
N
(
k
)
=
(
N
2
+
1
k
+
1
)
(
N
2
−
1
k
)
 as in Eq. (5).
S.-A. Cheong and C. L. Henley, Phys. Rev. B 80, 165124 (2009).
A. Russomanno, A. Silva, and G. E. Santoro, Phys. Rev. Lett. 109, 257201 (2012).
P. Weinberg, M. Bukov, L. D&apos;Alessio, A. Polkovnikov, S. Vajna, and M. Kolodrubetz, Phys. Rep. 688, 1 (2017).
J. Berges, S. Borsányi, and C. Wetterich, Phys. Rev. Lett. 93, 142002 (2004).
C. Kollath, A. M. Läuchli, and E. Altman, Phys. Rev. Lett. 98, 180601 (2007).
M. Eckstein, M. Kollar, and P. Werner, Phys. Rev. Lett. 103, 056403 (2009).
M. Moeckel and S. Kehrein, New J. Phys. 12, 055016 (2010).
L. Mathey and A. Polkovnikov, Phys. Rev. A 81, 033605 (2010).
R. Barnett, A. Polkovnikov, and M. Vengalattore, Phys. Rev. A 84, 023606 (2011).
M. Gring, M. Kuhnert, T. Langen, T. Kitagawa, B. Rauer, M. Schreitl, I. Mazets, D. A. Smith, E. Demler, and J. Schmiedmayer, Science 337, 1318 (2012).
F. Peronaci, M. Schiró, and O. Parcollet, , , and , arXiv:1711.07889.
M. Bukov, L. D&apos;Alessio, and A. Polkovnikov, Adv. Phys. 64, 139 (2015).
N. Goldman and J. Dalibard, Phys. Rev. X 4, 031027 (2014).
S. Rahav, I. Gilary, and S. Fishman, Phys. Rev. A 68, 013820 (2003).
V. Oganesyan and D. A. Huse, Phys. Rev. B 75, 155111 (2007).
We have numerically confirmed that the 
˜
H
±
1
 corrections to 
H
[
0
]
eff
 play a subleading role in all of the analyses in this paper.</item>
    </referencetext>
    <rights>No commercial reproduction, distribution, display or performance rights in this work are provided.</rights>
    <funders>
      <item>
        <agency>Department of Energy (DOE)</agency>
        <grant_number>DEAC02-05CH11231</grant_number>
      </item>
      <item>
        <agency>NSF</agency>
        <grant_number>DMR-1410435</grant_number>
      </item>
      <item>
        <agency>Institute of Quantum Information and Matter (IQIM)</agency>
      </item>
      <item>
        <agency>Gordon and Betty Moore Foundation</agency>
      </item>
      <item>
        <agency>David and Lucile Packard Foundation</agency>
      </item>
      <item>
        <agency>Army Research Office (ARO)</agency>
        <grant_number>W911NF-16-1-0361</grant_number>
      </item>
      <item>
        <agency>NSF Graduate Research Fellowship</agency>
      </item>
      <item>
        <agency>National Research Council of Canada</agency>
      </item>
      <item>
        <agency>Army Research Laboratory</agency>
      </item>
      <item>
        <agency>Air Force Office of Scientific Research (AFOSR)</agency>
      </item>
    </funders>
    <collection>CaltechAUTHORS</collection>
    <reviewer>George Porter</reviewer>
    <local_group>
      <item>IQIM</item>
      <item>Institute for Quantum Information and Matter</item>
    </local_group>
  </eprint>
</eprints>
`)
	// Test parsing
	rec := new(EPrints)
	err := xml.Unmarshal(src, &rec)
	if err != nil {
		t.Errorf("Can't parse eprint id 84590.xml, %s", err)
		t.FailNow()
	}
	if len(rec.EPrint) != 1 {
		t.Errorf("Wrongly structured of EPrint, %+v", rec)
	}
	eprint := rec.EPrint[0]
	if eprint.EPrintID != 84590 {
		t.Errorf("Expected eprint id of 84590, got %d", eprint.EPrintID)
	}
}

func TestMultipleEPrintResponse(t *testing.T) {
	// Simulate URL response for https://authors.library.caltech.edu/cgi/exportview?values=Accelerated_Strategic_Computing_Initiative&format=XML
	src := []byte(`<?xml version='1.0' encoding='utf-8'?>
<eprints xmlns='http://eprints.org/ep2/data/2.0'>
  <eprint id='https://authors.library.caltech.edu/id/eprint/25792'>
    <eprintid>25792</eprintid>
    <rev_number>9</rev_number>
    <documents>
      <document id='https://authors.library.caltech.edu/id/document/30675'>
        <docid>30675</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/358728'>
            <fileid>358728</fileid>
            <datasetid>document</datasetid>
            <objectid>30675</objectid>
            <filename>cit-asci-tr076.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <filesize>30949327</filesize>
            <mtime>2012-12-26 13:48:26</mtime>
            <url>https://authors.library.caltech.edu/25792/1/cit-asci-tr076.pdf</url>
          </file>
        </files>
        <eprintid>25792</eprintid>
        <pos>1</pos>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <security>public</security>
        <license>other</license>
        <main>cit-asci-tr076.pdf</main>
        <relation>
          <item>
            <type>http://eprints.org/relation/hasVolatileVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66129</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/haspreviewThumbnailVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66129</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/hasVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66129</uri>
          </item>
        </relation>
      </document>
    </documents>
    <eprint_status>archive</eprint_status>
    <userid>4298</userid>
    <dir>disk0/00/02/57/92</dir>
    <datestamp>2003-04-16</datestamp>
    <lastmod>2016-09-22 22:35:45</lastmod>
    <status_changed>2011-10-04 20:01:54</status_changed>
    <type>monograph</type>
    <metadata_visibility>show</metadata_visibility>
    <item_issues_count>0</item_issues_count>
    <creators>
      <item>
        <name>
          <family>Aivazis</family>
          <given>Michael</given>
        </name>
        <id>Aivazis-M</id>
      </item>
      <item>
        <name>
          <family>Goddard</family>
          <given>Bill</given>
        </name>
        <id>Goddard-W-A-III</id>
      </item>
      <item>
        <name>
          <family>Meiron</family>
          <given>Dan</given>
        </name>
        <id>Meiron-D-I</id>
      </item>
      <item>
        <name>
          <family>Ortiz</family>
          <given>Michael</given>
        </name>
        <id>Ortiz-M</id>
      </item>
      <item>
        <name>
          <family>Pool</family>
          <given>James C. T.</given>
        </name>
        <id>Pool-J-C-T</id>
      </item>
      <item>
        <name>
          <family>Shepherd</family>
          <given>Joe</given>
        </name>
        <id>Shepherd-J-E</id>
        <orcid>0000-0003-3181-9310</orcid>
      </item>
    </creators>
    <title>ASCI Alliance Center for Simulation of Dynamic Response in Materials FY 2000 Annual Report</title>
    <ispublished>unpub</ispublished>
    <full_text_status>public</full_text_status>
    <monograph_type>technical_report</monograph_type>
    <abstract>Introduction: 
This annual report describes research accomplishments for FY 00 of the Center for
Simulation of Dynamic Response of Materials. The Center is constructing a virtual
shock physics facility in which the full three dimensional response of a variety of target
materials can be computed for a wide range of compressive, tensional, and shear
loadings, including those produced by detonation of energetic materials. The goals
are to facilitate computation of a variety of experiments in which strong shock and
detonation waves are made to impinge on targets consisting of various combinations
of materials, compute the subsequent dynamic response of the target materials, and
validate these computations against experimental data.</abstract>
    <date>2000-01-01</date>
    <date_type>published</date_type>
    <id_number>CaltechASCI:2000.076</id_number>
    <institution>California Institute of Technology</institution>
    <department>ASCI Center for Simulation of Dynamic Response in Materials</department>
    <refereed>FALSE</refereed>
    <official_url>http://resolver.caltech.edu/CaltechASCI:2000.076</official_url>
    <related_url>
      <item>
        <url>http://www.cacr.caltech.edu/ASAP/onlineresources/publications/</url>
        <type>pub</type>
      </item>
    </related_url>
    <rights>You are granted permission for individual, educational, research and non-commercial reproduction, distribution, display and performance of this work in any format.</rights>
    <collection>CaltechASCI</collection>
    <local_group>
      <item>Accelerated Strategic Computing Initiative</item>
      <item>GALCIT</item>
    </local_group>
  </eprint>
  <eprint id='https://authors.library.caltech.edu/id/eprint/25791'>
    <eprintid>25791</eprintid>
    <rev_number>8</rev_number>
    <documents>
      <document id='https://authors.library.caltech.edu/id/document/30674'>
        <docid>30674</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/358721'>
            <fileid>358721</fileid>
            <datasetid>document</datasetid>
            <objectid>30674</objectid>
            <filename>cit-asci-tr033.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <filesize>24245514</filesize>
            <mtime>2012-12-26 13:48:25</mtime>
            <url>https://authors.library.caltech.edu/25791/1/cit-asci-tr033.pdf</url>
          </file>
        </files>
        <eprintid>25791</eprintid>
        <pos>1</pos>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <security>public</security>
        <license>other</license>
        <main>cit-asci-tr033.pdf</main>
        <relation>
          <item>
            <type>http://eprints.org/relation/hasVolatileVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66128</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/haspreviewThumbnailVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66128</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/hasVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66128</uri>
          </item>
        </relation>
      </document>
    </documents>
    <eprint_status>archive</eprint_status>
    <userid>4298</userid>
    <dir>disk0/00/02/57/91</dir>
    <datestamp>2001-07-16</datestamp>
    <lastmod>2016-09-22 22:34:45</lastmod>
    <status_changed>2011-10-04 20:01:52</status_changed>
    <type>monograph</type>
    <metadata_visibility>show</metadata_visibility>
    <item_issues_count>0</item_issues_count>
    <creators>
      <item>
        <name>
          <family>Aivazis</family>
          <given>Michael</given>
        </name>
        <id>Aivazis-M</id>
      </item>
      <item>
        <name>
          <family>Goddard</family>
          <given>Bill</given>
        </name>
        <id>Goddard-W-A-III</id>
      </item>
      <item>
        <name>
          <family>Meiron</family>
          <given>Dan</given>
        </name>
        <id>Meiron-D-I</id>
      </item>
      <item>
        <name>
          <family>Ortiz</family>
          <given>Michael</given>
        </name>
        <id>Ortiz-M</id>
      </item>
      <item>
        <name>
          <family>Pool</family>
          <given>James C. T.</given>
        </name>
        <id>Pool-J-C-T</id>
      </item>
      <item>
        <name>
          <family>Shepherd</family>
          <given>Joe</given>
        </name>
        <id>Shepherd-J-E</id>
        <orcid>0000-0003-3181-9310</orcid>
      </item>
    </creators>
    <title>The 1999 Center for Simulation of Dynamic Response in Materials Annual Technical Report</title>
    <ispublished>unpub</ispublished>
    <full_text_status>public</full_text_status>
    <monograph_type>technical_report</monograph_type>
    <abstract>Introduction: 
This annual report describes research accomplishments for FY 99 of the Center
for Simulation of Dynamic Response of Materials. The Center is constructing a
virtual shock physics facility in which the full three dimensional response of a
variety of target materials can be computed for a wide range of compressive, ten-
sional, and shear loadings, including those produced by detonation of energetic
materials. The goals are to facilitate computation of a variety of experiments
in which strong shock and detonation waves are made to impinge on targets
consisting of various combinations of materials, compute the subsequent dy-
namic response of the target materials, and validate these computations against
experimental data.</abstract>
    <date>1999-01-01</date>
    <date_type>published</date_type>
    <id_number>CaltechASCI:1999.033</id_number>
    <institution>California Institute of Technology</institution>
    <department>ASCI Center for Simulation of Dynamic Response of Materials</department>
    <refereed>FALSE</refereed>
    <official_url>http://resolver.caltech.edu/CaltechASCI:1999.033</official_url>
    <related_url>
      <item>
        <url>http://www.cacr.caltech.edu/ASAP/onlineresources/publications/</url>
        <type>pub</type>
      </item>
    </related_url>
    <rights>You are granted permission for individual, educational, research and non-commercial reproduction, distribution, display and performance of this work in any format.</rights>
    <collection>CaltechASCI</collection>
    <local_group>
      <item>Accelerated Strategic Computing Initiative</item>
      <item>GALCIT</item>
    </local_group>
  </eprint>
  <eprint id='https://authors.library.caltech.edu/id/eprint/25790'>
    <eprintid>25790</eprintid>
    <rev_number>8</rev_number>
    <documents>
      <document id='https://authors.library.caltech.edu/id/document/30673'>
        <docid>30673</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/358715'>
            <fileid>358715</fileid>
            <datasetid>document</datasetid>
            <objectid>30673</objectid>
            <filename>cit-asci-tr032.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <filesize>9971155</filesize>
            <mtime>2012-12-26 13:48:24</mtime>
            <url>https://authors.library.caltech.edu/25790/1/cit-asci-tr032.pdf</url>
          </file>
        </files>
        <eprintid>25790</eprintid>
        <pos>1</pos>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <security>public</security>
        <license>other</license>
        <main>cit-asci-tr032.pdf</main>
        <relation>
          <item>
            <type>http://eprints.org/relation/hasVolatileVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66127</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/haspreviewThumbnailVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66127</uri>
          </item>
          <item>
            <type>http://eprints.org/relation/hasVersion</type>
            <uri>https://authors.library.caltech.edu/id/document/66127</uri>
          </item>
        </relation>
      </document>
    </documents>
    <eprint_status>archive</eprint_status>
    <userid>4298</userid>
    <dir>disk0/00/02/57/90</dir>
    <datestamp>2001-07-16</datestamp>
    <lastmod>2016-09-22 22:34:32</lastmod>
    <status_changed>2011-10-04 20:01:50</status_changed>
    <type>monograph</type>
    <metadata_visibility>show</metadata_visibility>
    <item_issues_count>0</item_issues_count>
    <creators>
      <item>
        <name>
          <family>Goddard</family>
          <given>W. A.</given>
        </name>
        <id>Goddard-W-A-III</id>
      </item>
      <item>
        <name>
          <family>Meiron</family>
          <given>D. I.</given>
        </name>
        <id>Meiron-D-I</id>
      </item>
      <item>
        <name>
          <family>Ortiz</family>
          <given>M.</given>
        </name>
        <id>Ortiz-M</id>
      </item>
      <item>
        <name>
          <family>Shepherd</family>
          <given>J. E.</given>
        </name>
        <id>Shepherd-J-E</id>
        <orcid>0000-0003-3181-9310</orcid>
      </item>
    </creators>
    <title>The 1998 Center for Simulation of Dynamic Response in Materials Annual Technical Report</title>
    <ispublished>unpub</ispublished>
    <full_text_status>public</full_text_status>
    <monograph_type>technical_report</monograph_type>
    <abstract>Introduction: 
This annual report describes research accomplishments for FY 98 of the Center for Simulation
of Dynamic Response of Materials. The Center is constructing a virtual shock physics facility
in which the full three dimensional response of a variety of target materials can be computed
for a wide range of compressive, tensional, and shear loadings, including those produced by
detonation of energetic materials. The goals are to facilitate computation of a variety of
experiments in which strong shock and detonation waves are made to impinge on targets
consisting of various combinations of materials, compute the subsequent dynamic response
of the target materials, and validate these computations against experimental data.</abstract>
    <date>1998-01-01</date>
    <date_type>published</date_type>
    <id_number>CaltechASCI:1998.032</id_number>
    <institution>California Institute of Technology</institution>
    <department>ASCI Center for Simulation of Dynamic Response of Materials</department>
    <refereed>FALSE</refereed>
    <official_url>http://resolver.caltech.edu/CaltechASCI:1998.032</official_url>
    <related_url>
      <item>
        <url>http://www.cacr.caltech.edu/ASAP/onlineresources/publications/</url>
        <type>pub</type>
      </item>
    </related_url>
    <rights>You are granted permission for individual, educational, research and non-commercial reproduction, distribution, display and performance of this work in any format.</rights>
    <collection>CaltechASCI</collection>
    <local_group>
      <item>Accelerated Strategic Computing Initiative</item>
      <item>GALCIT</item>
    </local_group>
  </eprint>
</eprints>
`)

	records := new(EPrints)
	err := xml.Unmarshal(src, records)
	if err != nil {
		t.Errorf("Couldn't unmarshal multi-record response, %s", err)
		t.FailNow()
	}
	if len(records.EPrint) != 3 {
		t.Errorf("Expected 3 records, got %d", records.EPrint)
	}
}

func TestGetEPrint(t *testing.T) {
	getURL := os.Getenv("EP_EPRINT_URL")
	testKey := os.Getenv("EP_TEST_KEY")
	if getURL == "" || testKey == "" {
		t.Log("Skipping TestGetEPrint(), environment not setup")
		t.SkipNow()
	}
	t.Log("Runnning TestGetEPrint() with %s", getURL)
	restPath := "/rest/eprint/" + testKey + ".xml"
	u, _ := url.Parse(getURL + restPath)
	records := new(EPrints)
	records, xmlSrc, err := GetEPrints(u.String(), rc.AuthNone, "", "")
	if err != nil {
		t.Errorf("can't get %s, %s", u.String(), err)
	}
	if len(xmlSrc) == 0 {
		t.Errorf("Expected some XML data from %s", u.String())
	}
	if len(records.EPrint) == 0 {
		t.Errorf("Expected a populated record for %s, got %+v", testKey)
	}
	for i, rec := range records.EPrint {
		if rec.ID == "0" {
			t.Errorf("Expected a populated record key (%d) for %s, got %+v", i, testKey, rec)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	src := []byte(`
<?xml version='1.0' encoding='utf-8'?>
<eprints xmlns='http://eprints.org/ep2/data/2.0'>
  <eprint id='https://authors.library.caltech.edu/id/eprint/58360'>
    <eprintid>58360</eprintid>
    <rev_number>16</rev_number>
    <documents>
      <document id='https://authors.library.caltech.edu/id/document/218414'>
        <docid>218414</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/931167'>
            <fileid>931167</fileid>
            <datasetid>document</datasetid>
            <objectid>218414</objectid>
            <filename>0004-637X_805_1_8.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <hash>a4ad22c244df31a919bff6d7c8667fe9</hash>
            <hash_type>MD5</hash_type>
            <filesize>1209530</filesize>
            <mtime>2015-06-19 15:28:54</mtime>
            <url>https://authors.library.caltech.edu/58360/1/0004-637X_805_1_8.pdf</url>
          </file>
        </files>
        <eprintid>58360</eprintid>
        <pos>1</pos>
        <placement>1</placement>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <language>en</language>
        <security>public</security>
        <license>other</license>
        <main>0004-637X_805_1_8.pdf</main>
        <content>published</content>
      </document>
      <document id='https://authors.library.caltech.edu/id/document/218415'>
        <docid>218415</docid>
        <rev_number>2</rev_number>
        <files>
          <file id='https://authors.library.caltech.edu/id/file/931173'>
            <fileid>931173</fileid>
            <datasetid>document</datasetid>
            <objectid>218415</objectid>
            <filename>1501.04107v2.pdf</filename>
            <mime_type>application/pdf</mime_type>
            <hash>c475d5def259d6f04e54aa1dd363de1d</hash>
            <hash_type>MD5</hash_type>
            <filesize>1251916</filesize>
            <mtime>2015-06-19 15:37:19</mtime>
            <url>https://authors.library.caltech.edu/58360/2/1501.04107v2.pdf</url>
          </file>
        </files>
        <eprintid>58360</eprintid>
        <pos>2</pos>
        <placement>2</placement>
        <mime_type>application/pdf</mime_type>
        <format>application/pdf</format>
        <language>en</language>
        <security>public</security>
        <license>other</license>
        <main>1501.04107v2.pdf</main>
        <content>submitted</content>
      </document>
    </documents>
    <eprint_status>archive</eprint_status>
    <userid>772</userid>
    <dir>disk0/00/05/83/60</dir>
    <datestamp>2015-06-19 16:32:16</datestamp>
    <lastmod>2017-11-08 22:52:49</lastmod>
    <status_changed>2015-06-19 16:32:16</status_changed>
    <type>article</type>
    <metadata_visibility>show</metadata_visibility>
    <creators>
      <item>
        <name>
          <family>Zhu</family>
          <given>Wei</given>
        </name>
        <id>Zhu-Wei</id>
        <orcid>0000-0003-4027-4711</orcid>
      </item>
      <item>
        <name>
          <family>Udalski</family>
          <given>A.</given>
        </name>
        <id>Udalski-A</id>
      </item>
      <item>
        <name>
          <family>Gould</family>
          <given>A.</given>
        </name>
        <id>Gould-A</id>
      </item>
      <item>
        <name>
          <family>Dominik</family>
          <given>M.</given>
        </name>
        <id>Dominik-M</id>
      </item>
      <item>
        <name>
          <family>Bozza</family>
          <given>V.</given>
        </name>
        <id>Bozza-V</id>
      </item>
      <item>
        <name>
          <family>Han</family>
          <given>C.</given>
        </name>
        <id>Han-C</id>
      </item>
      <item>
        <name>
          <family>Yee</family>
          <given>J. C.</given>
        </name>
        <id>Yee-Jennifer-C</id>
      </item>
      <item>
        <name>
          <family>Calchi Novati</family>
          <given>S.</given>
        </name>
        <id>Calchi-Novati-S</id>
        <orcid>0000-0002-7669-1069</orcid>
      </item>
      <item>
        <name>
          <family>Beichman</family>
          <given>C. A.</given>
        </name>
        <id>Beichman-C-A</id>
      </item>
      <item>
        <name>
          <family>Carey</family>
          <given>S.</given>
        </name>
        <id>Carey-S-J</id>
        <orcid>0000-0002-0221-6871</orcid>
      </item>
      <item>
        <name>
          <family>Poleski</family>
          <given>R.</given>
        </name>
        <id>Poleski-R</id>
      </item>
      <item>
        <name>
          <family>Skowron</family>
          <given>J.</given>
        </name>
        <id>Skowron-J</id>
      </item>
      <item>
        <name>
          <family>Kozłowski</family>
          <given>S.</given>
        </name>
        <id>Kozłowski-S</id>
      </item>
      <item>
        <name>
          <family>Mrόz</family>
          <given>P.</given>
        </name>
        <id>Mrόz-P</id>
      </item>
      <item>
        <name>
          <family>Pietrukowicz</family>
          <given>P.</given>
        </name>
        <id>Pietrukowicz-P</id>
      </item>
      <item>
        <name>
          <family>Pietrzyński</family>
          <given>G.</given>
        </name>
        <id>Pietrzyński-G</id>
      </item>
      <item>
        <name>
          <family>Szymański</family>
          <given>M. K.</given>
        </name>
        <id>Szymański-M-K</id>
      </item>
      <item>
        <name>
          <family>Soszyński</family>
          <given>I.</given>
        </name>
        <id>Soszyński-I</id>
      </item>
      <item>
        <name>
          <family>Ulaczyk</family>
          <given>K.</given>
        </name>
        <id>Ulaczyk-K</id>
      </item>
      <item>
        <name>
          <family>Wyrzykowski</family>
          <given>Ł.</given>
        </name>
        <id>Wyrzykowski-Ł</id>
      </item>
    </creators>
    <corp_creators>
      <item>
        <name>OGLE Collaboration</name>
        <id>OGLE-Collaboration</id>
      </item>
      <item>
        <name>µFUN Collaboration</name>
        <id>µFUN-Collaboration</id>
      </item>
    </corp_creators>
    <title>Spitzer as a Microlens Parallax Satellite: Mass and Distance Measurements of the Binary Lens System OGLE-2014-BLG-1050L</title>
    <ispublished>pub</ispublished>
    <full_text_status>public</full_text_status>
    <keywords>binaries: general; gravitational lensing: micro</keywords>
    <note>© 2015 American Astronomical Society. 

Received 2015 January 20; accepted 2015 March 12; published 2015 May 13. 

Work by W.Z., A.G., and B.S.G. was supported by NSF grant AST 1103471. Work by J.C.Y., A.G., and S.C. was supported by JPL grant 1500811. A.G., B.S.G., and R.W.P. were supported by NASA grant NNX12AB99G. Work by C.H. was supported by the Creative Research Initiative Program (2009-0081561) of the National Research Foundation of Korea. Work by J.C.Y. was performed under contract with the California Institute of Technology (Caltech)/Jet Propulsion Laboratory (JPL) funded by NASA through the Sagan Fellowship Program executed by the NASA Exoplanet Science Institute. Work by C.A.B. was carried out in part at the Jet Propulsion Laboratory (JPL), California Institute of Technology, under a contract with the National Aeronautics and Space Administration. The OGLE project has received funding from the European Research Council under the European Community’s Seventh Framework Programme (FP7/2007–2013)/ERC grant agreement no. 246678 to AU. This work is based in part on observations made with the Spitzer Space Telescope, which is operated by the Jet Propulsion Laboratory, California Institute of Technology, under a contract with NASA.</note>
    <abstract>We report the first mass and distance measurements of a caustic-crossing binary system OGLE-2014-BLG-1050 L using the space-based microlens parallax method. Spitzer captured the second caustic crossing of the event, which occurred ~10 days before that seen from Earth. Due to the coincidence that the source-lens relative motion was almost parallel to the direction of the binary-lens axis, the fourfold degeneracy, which was known before only to occur in single-lens events, persists in this case, leading to either a lower-mass (0.2 and 0.07 M_☉) binary at ~1.1 kpc or a higher-mass (0.9 and 0.35 M_☉) binary at ~3.5 kpc. However, the latter solution is strongly preferred for reasons including blending and lensing probability. OGLE-2014-BLG-1050 L demonstrates the power of microlens parallax in probing stellar and substellar binaries.</abstract>
    <date>2015-05-20</date>
    <date_type>published</date_type>
    <publication>Astrophysical Journal</publication>
    <volume>805</volume>
    <number>1</number>
    <publisher>American Astronomical Society</publisher>
    <pagerange>Art. No. 8</pagerange>
    <id_number>CaltechAUTHORS:20150619-082837114</id_number>
    <refereed>TRUE</refereed>
    <issn>0004-637X</issn>
    <official_url>http://resolver.caltech.edu/CaltechAUTHORS:20150619-082837114</official_url>
    <related_url>
      <item>
        <url>http://dx.doi.org/10.1088/0004-637X/805/1/8</url>
        <type>doi</type>
        <description>Article</description>
      </item>
      <item>
        <url>http://iopscience.iop.org/0004-637X/805/1/8/</url>
        <type>pub</type>
        <description>Article</description>
      </item>
      <item>
        <url>http://arxiv.org/abs/1501.04107</url>
        <type>arxiv</type>
        <description>Discussion Paper</description>
      </item>
    </related_url>
    <referencetext>
      <item>Albrow, M. D., Beaulieu, J.-P., Caldwell, J. A. R., et al. 1999, ApJ, 522, 1022
Alcock, C., Allsman, R. A., Alves, D., et al. 1995, ApJL, 454, L125
Alcock, C., Allsman, R. A., Alves, D. R., et al. 2001, Natur, 414, 617
Batista, V., Beaulieu, J.-P., Gould, A., et al. 2014, ApJ, 780, 54
Batista, V., Gould, A., Dieters, S., et al. 2011, A&amp;A, 529, A102
Bennett, D. P., Anderson, J., Bond, I. A., Udalski, A., &amp; Gould, A. 2006,
ApJL, 647, L171
Bensby, T., Yee, J. C., Feltzing, S., et al. 2013, A&amp;A, 549, A147
Bessell, M. S., &amp; Brett, J. M. 1988, PASP, 100, 1134
Bozza, V. 2010, MNRAS, 408, 2188
Calchi Novati, S., Gould, A., Udalski, A., et al. 2014, arXiv:1411.7378
Cassan, A. 2008, A&amp;A, 491, 587
Cassan, A., Horne, K., Kains, N., Tsapras, Y., &amp; Browne, P. 2010, A&amp;A,
515, A52
Choi, J.-Y., Han, C., Udalski, A., et al. 2013, ApJ, 768, 129
Dominik, M. 1998, A&amp;A, 333, L79
Dong, S., DePoy, D. L., Gaudi, B. S., et al. 2006, ApJ, 642, 842
Dong, S., Gould, A., Udalski, A., et al. 2009, ApJ, 695, 970
Dong, S., Udalski, A., Gould, A., et al. 2007, ApJ, 664, 862
Duchêne, G., &amp; Kraus, A. 2013, ARA&amp;A, 51, 269
Gaudi, B. S., Bennett, D. P., Udalski, A., et al. 2008, Sci, 319, 927
Gould, A. 1992, ApJ, 392, 442
Gould, A. 1994, ApJL, 421, L75
Gould, A. 1997, ApJ, 480, 188
Gould, A. 2004, ApJ, 606, 319
Gould, A. 2008, ApJ, 681, 1593
Gould, A. 2013, ApJL, 763, L35
Gould, A. 2014, JKAS, 47, 215
Gould, A., Carey, S., &amp; Yee, J. 2014, sptz prop, 11006
Gould, A., &amp; Gaucherel, C. 1997, ApJ, 477, 580
Gould, A., &amp; Horne, K. 2013, ApJL, 779, L28
Gould, A., Udalski, A., Monard, B., et al. 2009, ApJL, 698, L147
Gould, A., &amp; Yee, J. C. 2012, ApJL, 755, L17
Gould, A., &amp; Yee, J. C. 2014, ApJ, 784, 64
Graff, D. S., &amp; Gould, A. 2002, ApJ, 580, 253
Hardy, S. J., &amp; Walker, M. A. 1995, MNRAS, 276, L79
Honma, M. 1999, ApJL, 517, L35
Jiang, G., DePoy, D. L., Gal-Yam, A., et al. 2004, ApJ, 617, 1307
Jung, Y. K., Udalski, A., Sumi, T., et al. 2015, ApJ, 798, 123
Kains, N., Cassan, A., Horne, K., et al. 2009, MNRAS, 395, 787
Kervella, P., Thévenin, F., di Folco, E., &amp; Ségransan, D. 2004, A&amp;A, 426, 297
Nataf, D. M., Gould, A., Fouqué, P., et al. 2013, ApJ, 769, 88
Park, H., Udalski, A., Han, C., et al. 2013, ApJ, 778, 134
Pejcha, O., &amp; Heyrovský, D. 2009, ApJ, 690, 1772
Penny, M. T., Kerins, E., Rattenbury, N., et al. 2013, MNRAS, 434, 2
Poindexter, S., Afonso, C., Bennett, D. P., et al. 2005, ApJ, 633, 914
Ramírez, I., Michel, R., Sefako, R., et al. 2012, ApJ, 752, 5
Refsdal, S. 1966, MNRAS, 134, 315
Spergel, D., Gehrels, N., Breckinridge, J., et al. 2013, arXiv:1305.5422
Skowron, J., Udalski, A., Gould, A., et al. 2011, ApJ, 738, 87
Udalski, A. 2003, AcA, 53, 291
Udalski, A., Yee, J. C., Gould, A., et al. 2015, ApJ, 799, 237
Walker, M. A. 1995, ApJ, 453, 37
Yee, J. C., Hung, L.-W., Bond, I. A., et al. 2013, ApJ, 769, 77
Yee, J. C., Udalski, A., Calchi Novati, S., et al. 2015, ApJ, 802, 76
Yee, J. C., Udalski, A., Sumi, T., et al. 2009, ApJ, 703, 2082
Yoo, J., DePoy, D. L., Gal-Yam, A., et al. 2004, ApJ, 603, 139
Zhu, W., Penny, M., Mao, S., Gould, A., &amp; Gendron, R. 2014, ApJ, 788, 73</item>
    </referencetext>
    <rights>No commercial reproduction, distribution, display or performance rights in this work are provided.</rights>
    <official_cit>Spitzer as a Microlens Parallax Satellite: Mass and Distance Measurements of the Binary Lens System OGLE-2014-BLG-1050L
Wei Zhu (祝伟) et al. 2015 ApJ 805 8</official_cit>
    <funders>
      <item>
        <agency>NSF</agency>
        <grant_number>AST-1103471</grant_number>
      </item>
      <item>
        <agency>JPL</agency>
        <grant_number>1500811</grant_number>
      </item>
      <item>
        <agency>NASA</agency>
        <grant_number>NNX12AB99G</grant_number>
      </item>
      <item>
        <agency>National Research Foundation of Korea</agency>
        <grant_number>2009-0081561</grant_number>
      </item>
      <item>
        <agency>NASA Sagan Fellowship</agency>
      </item>
      <item>
        <agency>NASA/JPL/Caltech</agency>
      </item>
      <item>
        <agency>European Research Council (ERC)</agency>
        <grant_number>246678</grant_number>
      </item>
    </funders>
    <collection>CaltechAUTHORS</collection>
    <reviewer>George Porter</reviewer>
    <local_group>
      <item>Infrared Processing and Analysis Center (IPAC)</item>
    </local_group>
  </eprint>
</eprints>
`)
	records := new(EPrints)
	err := xml.Unmarshal(src, &records)
	if err != nil {
		t.Errorf("unmashal failed, %s", err)
		t.FailNow()
	}
	if len(records.EPrint) != 1 {
		t.Errorf("Should have found at least one eprint")
		t.FailNow()
	}
	record := records.EPrint[0]
	if record.Publication != "Astrophysical Journal" {
		t.Errorf("expected %q, got %q", "Astrophysical Journal", record.Publication)
	}
	if record.Volume != "805" {
		t.Errorf("expected %q, got %q", "805", record.Volume)
	}
	if record.Number != "1" {
		t.Errorf("expected %q, got %q", "1", record.Number)
	}
}
