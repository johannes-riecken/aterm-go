record Typ where
  constructor MkType
  name : String
  -- obj : Maybe Object

record Field where
  names : List Ident
  type : Expr
record FieldList where
  constructor MkFieldList
  list : List Field

record FuncType where
  constructor MkFuncType
  params : FieldList

record FuncDecl where
  constructor MkFuncDecl
  name : Ident
  type : FuncType

record Object where
  constructor FunObj
  name : String
  decl : FuncDecl
  -- constructor BarObj

record Ident where
  constructor MkIdent
  name : String
  obj : Maybe Object

record Decl where
  constructor MkDecl
  name : Ident

record File where
  constructor MkFile
  name : Ident
  decls : List Decl

-- record Object
--   constructor MkObject
--   name : String



