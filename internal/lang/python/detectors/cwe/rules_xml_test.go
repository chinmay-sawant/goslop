package cwe_test

import "testing"

func TestCWEXMLHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		id   string
		src  string
		hit  bool
	}{
		{
			name: "CWE-611 lxml entity resolution", id: "CWE-611", hit: true,
			src: "from lxml import etree\n\ndef parse(data):\n    parser = etree.XMLParser(resolve_entities=True)\n    return etree.fromstring(data, parser)\n",
		},
		{
			name: "CWE-611 entity resolution disabled", id: "CWE-611", hit: false,
			src: "from lxml import etree\n\ndef parse(data):\n    parser = etree.XMLParser(resolve_entities=False, no_network=True)\n    return etree.fromstring(data, parser)\n",
		},
		{
			name: "CWE-611 XMLParser in comment", id: "CWE-611", hit: false,
			src: "# etree.XMLParser(resolve_entities=True)\nvalue = 'etree.XMLParser(resolve_entities=True)'\n",
		},
		{
			name: "CWE-776 DTD entity expansion", id: "CWE-776", hit: true,
			src: "from lxml import etree\n\ndef parse(data):\n    parser = etree.XMLParser(load_dtd=True, resolve_entities=True, huge_tree=True)\n    return etree.fromstring(data, parser)\n",
		},
		{
			name: "CWE-776 DTD without entity resolution", id: "CWE-776", hit: false,
			src: "from lxml import etree\n\ndef parse(data):\n    parser = etree.XMLParser(load_dtd=True, resolve_entities=False)\n    return etree.fromstring(data, parser)\n",
		},
		{
			name: "CWE-112 request XML without schema", id: "CWE-112", hit: true,
			src: "from lxml import etree\n\ndef parse(request):\n    return etree.fromstring(request.data)\n",
		},
		{
			name: "CWE-112 schema-aware parser", id: "CWE-112", hit: false,
			src: "from lxml import etree\n\ndef parse(request, schema):\n    parser = etree.XMLParser(schema=schema, resolve_entities=False)\n    return etree.fromstring(request.data, parser=parser)\n",
		},
		{
			name: "CWE-112 application-controlled XML", id: "CWE-112", hit: false,
			src: "from xml.etree import ElementTree\n\ndef parse():\n    return ElementTree.fromstring('<config/>')\n",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := hasRule(runScan(tc.src), tc.id)
			if got != tc.hit {
				t.Fatalf("%s hit=%v, want %v", tc.id, got, tc.hit)
			}
		})
	}
}
