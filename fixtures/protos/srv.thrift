
struct TestRep {
    1: string ID,
    2: binary Data,
}

service TSrv {
    TestRep Test(
        1: TestRep par,
    ),
}