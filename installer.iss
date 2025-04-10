[Setup]
AppName=Media Manager
AppVersion=1.0.0
DefaultDirName={pf}\MediaManager
DefaultGroupName=Media Manager
OutputDir=output
OutputBaseFilename=MediaManagerSetup
Compression=lzma
SolidCompression=yes

[Files]
Source: "dist\mediamanager.exe"; DestDir: "{app}"
Source: "dist\bin\*"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "dist\database\*"; DestDir: "{app}\database"; Flags: recursesubdirs
Source: "dist\video.html"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Media Manager"; Filename: "{app}\mediamanager.exe"
Name: "{commondesktop}\Media Manager"; Filename: "{app}\mediamanager.exe"

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "PATH"; ValueData: "{olddata};{app}\bin"; Check: NeedsAddPath(ExpandConstant('{app}\bin'))

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKLM, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PATH', OrigPath) then
  begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;
