using Go = import "/go.capnp";
$Go.package("mdbxsql");
$Go.import("mdbxsql/capnp");

@0x91a877877ad1a9c6;

struct User {
  id @0 :Int32;
  name @1 :Text;
  age @2 :Int32;
}
