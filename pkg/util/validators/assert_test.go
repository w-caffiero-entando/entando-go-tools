package validators

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
    assert.Equal(t, Assert_StrictId("STRICT ID", "object", "", "silent"), true)
    assert.Equal(t, Assert_StrictId("STRICT ID", "my_object", "", "silent"), false)
    assert.Equal(t, Assert_StrictId("STRICT ID", "my-object", "", "silent"), false)
    assert.Equal(t, Assert_StrictId("STRICT ID", "myObject", "", "silent"), true)
    assert.Equal(t, Assert_StrictId("STRICT ID", "MyObject", "", "silent"), false)
    assert.Equal(t, Assert_StrictId("STRICT ID", "my object", "", "silent"), false)
    assert.Equal(t, Assert_StrictId("STRICT ID", "my.Object", "", "silent"), false)

    assert.Equal(t, Assert_Id("ID", "object", "", "silent"), true)
    assert.Equal(t, Assert_Id("ID", "my_object", "", "silent"), true)
    assert.Equal(t, Assert_Id("ID", "my-object", "", "silent"), false)
    assert.Equal(t, Assert_Id("ID", "myObject", "", "silent"), true)
    assert.Equal(t, Assert_Id("ID", "MyObject", "", "silent"), false)
    assert.Equal(t, Assert_Id("ID", "my object", "", "silent"), false)
    assert.Equal(t, Assert_Id("ID", "my.Object", "", "silent"), false)

    assert.Equal(t, Assert_Id_Spc("SPC ID", "object", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "my_object", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "my-object", "", "silent"), false)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "myObject", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "MyObject", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "my object", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "My Object", "", "silent"), true)
    assert.Equal(t, Assert_Id_Spc("SPC ID", "my.Object", "", "silent"), false)

    assert.Equal(t, Assert_Ext_Ic_Id_Spc("IC ID SPC", "My Object", "", "silent"), true)
    assert.Equal(t, Assert_Ext_Ic_Id_Spc("IC ID SPC", "My-Object", "", "silent"), true)
    assert.Equal(t, Assert_Ext_Ic_Id_Spc("IC ID SPC", "My- Object", "", "silent"), true)

    assert.Equal(t, Assert_DN("DN", "domain", "", "silent"), true)
    assert.Equal(t, Assert_DN("DN", "domai@n", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "domain?", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "domain99", "", "silent"), true)
    assert.Equal(t, Assert_DN("DN", "my_domain", "", "silent"), true)
    assert.Equal(t, Assert_DN("DN", "my-domain", "", "silent"), true)
    assert.Equal(t, Assert_DN("DN", "mydomain.example.com", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "myDomain", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "MyDomain", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "my domain", "", "silent"), false)
    assert.Equal(t, Assert_DN("DN", "my.Domain", "", "silent"), false)

    assert.Equal(t, Assert_FDN("FDN", "domain", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "d@omain", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "domain?", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "domain99", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "my_domain", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "my-domain", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "mydomain.example.com", "", "silent"), true)
    assert.Equal(t, Assert_FDN("FDN", "my-d_omain.example.com", "", "silent"), true)
    assert.Equal(t, Assert_FDN("FDN", "myDomain.example.com", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "myDomain.example.com", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "MyDomain.example.com", "", "silent"), false)
    assert.Equal(t, Assert_FDN("FDN", "my domain.example.com", "", "silent"), false)

    assert.Equal(t, Assert_URL("URL", "domain", "", "silent"), false)
    assert.Equal(t, Assert_URL("URL", "http://example.com", "", "silent"), true)
    assert.Equal(t, Assert_URL("URL", "http://example.com?q=1#f", "", "silent"), true)

    assert.Equal(t, Assert_EMAIL("MAIL", "domain", "", "silent"), false)
    assert.Equal(t, Assert_EMAIL("MAIL", "domain@example.com", "", "silent"), true)
    assert.Equal(t, Assert_EMAIL("MAIL", "asb.s@example.com", "", "silent"), true)
    assert.Equal(t, Assert_EMAIL("MAIL", "asb.s@an.example.com", "", "silent"), true)
    assert.Equal(t, Assert_EMAIL("MAIL", ".s@example.com", "", "silent"), false)

    assert.Equal(t, Assert_IP("IP", "192.168.0.1", "", "silent"), true)
    assert.Equal(t, Assert_IP("IP", "", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "....", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", ".168.1.0", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "192.168..0", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "192.168.0.", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "192.168.0", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "192.168.0.a1", "", "silent"), false)
    assert.Equal(t, Assert_IP("IP", "192.168a.0.1", "", "silent"), false)

    assert.Equal(t, Assert_Ver("VER", "0.0.1", "", "silent"), true)
    assert.Equal(t, Assert_Ver("VER", "0.0.1-SNAPSHOT", "", "silent"), true)
    assert.Equal(t, Assert_Ver("VER", "0.0.1-SNAPSHOT-2", "", "silent"), true)
    assert.Equal(t, Assert_Ver("VER", "0.1-SNAPSHOT", "", "silent"), true)
    assert.Equal(t, Assert_Ver("VER", "1-SNAPSHOT", "", "silent"), true)
    assert.Equal(t, Assert_Ver("VER", "0.0.1+SNAPSHOT", "", "silent"), false)
    assert.Equal(t, Assert_Ver("VER", "0.0.1:2", "", "silent"), false)
    assert.Equal(t, Assert_Ver("VER", "0.0.1_3", "", "silent"), false)
    assert.Equal(t, Assert_Ver("VER", "v0.0.1", "", "silent"), true)

    assert.Equal(t, Assert_Giga("GIGA", "10G", "", "silent"), true)
    assert.Equal(t, Assert_Giga("GIGA", "10GG", "", "silent"), false)
    assert.Equal(t, Assert_Giga("GIGA", "10", "", "silent"), false)
}
