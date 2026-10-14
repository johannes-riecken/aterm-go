data File = MkFile { name :: Ident [Decl]

newtype FuncType = MkFuncType FieldList

newtype FieldList = MkFieldList [Field]

data Field = MkField { fieldName :: Ident, fieldType :: Ident }

-- MkIdent "foo" (Just (MkObject "var" "foo"))
data Ident = MkIdent String (Maybe Object)

-- MkObject "var" "foobar" 0 (TODO: meaning of 0?)
data Object = MkObject String String

data BlockStmt = MkBlockStmt [DeclStmt]

data DeclStmt = MkDeclStmt GenDecl

-- MkGenDecl "var" ...
data GenDecl = MkGenDecl String [ValueSpec]

data ValueSpec = MkValueSpec [Ident]

funcType :: FuncType
funcType = MkFuncType (MkFieldList [MkField (MkIdent "foo" (Just $ MkObject "var" "foo")) (MkIdent "int" Nothing)])

main :: IO ()
main = pure ()

File(Ident("main"),[FuncDecl(Ident("exampleFunc",Object("func","exampleFunc")),FuncType(FieldList([Field([Ident("foo",Object("var","foo"))],Ident("int")),Field([Ident("bar",Object("var","bar"))],ArrayType(Ident("int"))),Field([Ident("baz",Object("var","baz"))],Ident("int")),Field([Ident("quux",Object("var","quux"))],ArrayType(Ident("int")))])),

BlockStmt([
    DeclStmt(GenDecl("var",[ValueSpec([Ident("foobar",Object("var","foobar",0))],Ident("int"))])),
    AssignStmt([Ident("foobar"),Ident("foo"),Ident("quux")],"=",[CallExpr(Ident("len"),[Ident("bar")]),CallExpr(Ident("len"),[Ident("quux")]),CompositeLit(ArrayType(Ident("int")),[Ident("foo"),Ident("baz")],False)]),
    AssignStmt([Ident("_")],"=",[Ident("foobar")])])

)],Scope(),[],"")

