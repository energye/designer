unit Unit1;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Math,
  Forms, Controls, Graphics, Dialogs, StdCtrls, ExtCtrls;

type
  TForm1 = class(TForm)
    procedure FormCreate(Sender: TObject);
  private
    HostPanel: TPanel;
    ViewPanel: TPanel;
    ContentPanel: TPanel;
    HBar: TScrollBar;
    VBar: TScrollBar;

    // 平滑滚动目标值
    FTargetX: Integer;
    FTargetY: Integer;
    FAnimTimer: TTimer;

    // 中键拖动
    FDragging: Boolean;
    FDragMouse: TPoint;
    FDragScroll: TPoint;

    procedure ScrollChange(Sender: TObject);
    procedure HostResize(Sender: TObject);

    procedure UpdateScrollBars;
    procedure UpdateView;
    procedure BuildManyControls(ACount: Integer);

    procedure AnimateScroll(Sender: TObject);

    procedure ContentMouseDown(Sender: TObject; Button: TMouseButton;
      Shift: TShiftState; X, Y: Integer);
    procedure ContentMouseMove(Sender: TObject; Shift: TShiftState;
      X, Y: Integer);
    procedure ContentMouseUp(Sender: TObject; Button: TMouseButton;
      Shift: TShiftState; X, Y: Integer);

    procedure ZoomCanvas(Delta: Integer);
    procedure SetScrollTarget(AX, AY: Integer);

    procedure FormMouseWheel(Sender: TObject; Shift: TShiftState;
      WheelDelta: Integer; MousePos: TPoint; var Handled: Boolean);
  end;

var
  Form1: TForm1;

implementation

{$R *.lfm}

procedure TForm1.FormCreate(Sender: TObject);
begin
  Width := 1200;
  Height := 800;
  Caption := 'Custom ScrollBox';

  DoubleBuffered := True;

  // Host Panel (Container)
  HostPanel := TPanel.Create(Self);
  HostPanel.Parent := Self;
  HostPanel.Align := alClient;
  HostPanel.BevelOuter := bvNone;
  HostPanel.DoubleBuffered := True;
  HostPanel.OnResize := @HostResize;

  // View Panel (Viewport)
  ViewPanel := TPanel.Create(Self);
  ViewPanel.Parent := HostPanel;
  ViewPanel.BevelOuter := bvNone;
  ViewPanel.Caption := '';
  ViewPanel.DoubleBuffered := True;

  // Content Panel (Main Content Area)
  ContentPanel := TPanel.Create(Self);
  ContentPanel.Parent := ViewPanel;
  ContentPanel.Caption := '';
  ContentPanel.DoubleBuffered := True;
  ContentPanel.SetBounds(0, 0, 6000, 4000);

  ContentPanel.OnMouseDown := @ContentMouseDown;
  ContentPanel.OnMouseMove := @ContentMouseMove;
  ContentPanel.OnMouseUp := @ContentMouseUp;

  // Scrollbars
  HBar := TScrollBar.Create(Self);
  HBar.Parent := HostPanel;
  HBar.Kind := sbHorizontal;
  HBar.OnChange := @ScrollChange;

  VBar := TScrollBar.Create(Self);
  VBar.Parent := HostPanel;
  VBar.Kind := sbVertical;
  VBar.OnChange := @ScrollChange;

  // Animation Timer (Smooth Scroll)
  FAnimTimer := TTimer.Create(Self);
  FAnimTimer.Interval := 15;
  FAnimTimer.Enabled := False;
  FAnimTimer.OnTimer := @AnimateScroll;

  // Build Initial Controls
  BuildManyControls(3000);

  UpdateScrollBars;
  UpdateView;

  // Enable Mouse Wheel Support
  Self.OnMouseWheel := @FormMouseWheel;
  HostPanel.OnMouseWheel := @FormMouseWheel;
  ViewPanel.OnMouseWheel := @FormMouseWheel;
  ContentPanel.OnMouseWheel := @FormMouseWheel;
end;

procedure TForm1.BuildManyControls(ACount: Integer);
const
  COLS = 20;
  W = 120;
  H = 32;
  GAP = 8;
var
  i, Row, Col: Integer;
  Btn: TButton;
begin
  ContentPanel.DisableAlign;
  try
    for i := 0 to ACount - 1 do
    begin
      Row := i div COLS;
      Col := i mod COLS;

      Btn := TButton.Create(Self);
      Btn.Parent := ContentPanel;
      Btn.SetBounds(
        Col * (W + GAP) + 10,
        Row * (H + GAP) + 10,
        W, H
      );
      Btn.Caption := 'Button ' + IntToStr(i + 1);
    end;
  finally
    ContentPanel.EnableAlign;
  end;
end;

procedure TForm1.HostResize(Sender: TObject);
const
  SB = 17;
var
  NeedH, NeedV: Boolean;
begin
  NeedH := ContentPanel.Width > HostPanel.ClientWidth;
  NeedV := ContentPanel.Height > HostPanel.ClientHeight;

  HBar.Visible := NeedH;
  VBar.Visible := NeedV;

  ViewPanel.SetBounds(
    0, 0,
    HostPanel.ClientWidth - Ord(NeedV) * SB,
    HostPanel.ClientHeight - Ord(NeedH) * SB
  );

  HBar.SetBounds(0, ViewPanel.Height, ViewPanel.Width, SB);
  VBar.SetBounds(ViewPanel.Width, 0, SB, ViewPanel.Height);

  UpdateScrollBars;
  UpdateView;
end;

procedure TForm1.UpdateScrollBars;
begin
  // Horizontal Scrollbar
  HBar.Min := 0;
  HBar.PageSize := ViewPanel.ClientWidth;
  HBar.Max := ContentPanel.Width;
  HBar.SmallChange := 20;
  HBar.LargeChange := 100;

  // Vertical Scrollbar
  VBar.Min := 0;
  VBar.PageSize := ViewPanel.ClientHeight;
  VBar.Max := ContentPanel.Height;
  VBar.SmallChange := 20;
  VBar.LargeChange := 100;

  HBar.Position := EnsureRange(HBar.Position, 0,
    Max(0, HBar.Max - Integer(HBar.PageSize)));

  VBar.Position := EnsureRange(VBar.Position, 0,
    Max(0, VBar.Max - Integer(VBar.PageSize)));

  FTargetX := HBar.Position;
  FTargetY := VBar.Position;
end;

procedure TForm1.UpdateView;
begin
  ContentPanel.Left := -HBar.Position;
  ContentPanel.Top := -VBar.Position;
end;

procedure TForm1.ScrollChange(Sender: TObject);
begin
  FTargetX := HBar.Position;
  FTargetY := VBar.Position;
  UpdateView;
end;

procedure TForm1.SetScrollTarget(AX, AY: Integer);
begin
  FTargetX := EnsureRange(AX, 0,
    Max(0, HBar.Max - Integer(HBar.PageSize)));

  FTargetY := EnsureRange(AY, 0,
    Max(0, VBar.Max - Integer(VBar.PageSize)));

  FAnimTimer.Enabled := True;
end;

procedure TForm1.AnimateScroll(Sender: TObject);
var
  DX, DY: Integer;
begin
  DX := FTargetX - HBar.Position;
  DY := FTargetY - VBar.Position;

  if Abs(DX) <= 1 then
    HBar.Position := FTargetX
  else
    HBar.Position := HBar.Position + DX div 4;

  if Abs(DY) <= 1 then
    VBar.Position := FTargetY
  else
    VBar.Position := VBar.Position + DY div 4;

  UpdateView;

  if (HBar.Position = FTargetX) and
     (VBar.Position = FTargetY) then
    FAnimTimer.Enabled := False;
end;

procedure TForm1.ContentMouseDown(Sender: TObject; Button: TMouseButton;
  Shift: TShiftState; X, Y: Integer);
begin
  if Button = mbMiddle then
  begin
    FDragging := True;
    FDragMouse := Mouse.CursorPos;
    FDragScroll := Point(FTargetX, FTargetY);
    Screen.Cursor := crSizeAll;
  end;
end;

procedure TForm1.ContentMouseMove(Sender: TObject; Shift: TShiftState;
  X, Y: Integer);
var
  P: TPoint;
begin
  if not FDragging then Exit;

  P := Mouse.CursorPos;

  SetScrollTarget(
    FDragScroll.X - (P.X - FDragMouse.X),
    FDragScroll.Y - (P.Y - FDragMouse.Y)
  );
end;

procedure TForm1.ContentMouseUp(Sender: TObject; Button: TMouseButton;
  Shift: TShiftState; X, Y: Integer);
begin
  if Button = mbMiddle then
  begin
    FDragging := False;
    Screen.Cursor := crDefault;
  end;
end;

procedure TForm1.ZoomCanvas(Delta: Integer);
var
  Scale: Double;
begin
  if Delta > 0 then
    Scale := 1.1
  else
    Scale := 0.9;

  ContentPanel.Width := Round(ContentPanel.Width * Scale);
  ContentPanel.Height := Round(ContentPanel.Height * Scale);

  if ContentPanel.Width < 500 then ContentPanel.Width := 500;
  if ContentPanel.Height < 400 then ContentPanel.Height := 400;

  HostResize(nil);
end;

procedure TForm1.FormMouseWheel(Sender: TObject; Shift: TShiftState;
  WheelDelta: Integer; MousePos: TPoint; var Handled: Boolean);
var
  Step: Integer;
begin
  Step := 80;

  // Ctrl + Wheel = Zoom
  if ssCtrl in Shift then
  begin
    ZoomCanvas(WheelDelta);
    Handled := True;
    Exit;
  end;

  // Shift + Wheel = Horizontal Scroll
  if ssShift in Shift then
  begin
    // Horizontal scroll
    SetScrollTarget(FTargetX - Sign(WheelDelta) * Step, FTargetY);
  end
  else
  begin
    // Vertical scroll
    SetScrollTarget(FTargetX, FTargetY - Sign(WheelDelta) * Step);
  end;

  Handled := True;
end;

end.