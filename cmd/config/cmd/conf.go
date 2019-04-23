package main

import (
	"fmt"
	"strings"

	"git.parallelcoin.io/dev/9/cmd/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const menutitle = "ⓟ parallelcoin 9 configuration CLI"

func MainColor() tcell.Color {
	return tcell.NewRGBColor(64, 64, 64)
}

func DimColor() tcell.Color {
	return tcell.NewRGBColor(48, 48, 48)
}

func PrelightColor() tcell.Color {
	return tcell.NewRGBColor(32, 32, 32)
}

func TextColor() tcell.Color {
	return tcell.NewRGBColor(216, 216, 216)
}

func BackgroundColor() tcell.Color {
	return tcell.NewRGBColor(16, 16, 16)
}

func Run(args []string, tokens config.Tokens, app *config.App) int {
	// tapp pulls everything together to create the configuration interface
	tapp := tview.NewApplication()

	// titlebar tells the user what app they are using
	titlebar := tview.NewTextView().
		SetTextColor(TextColor()).
		SetText(menutitle)
	titlebar.Box.SetBackgroundColor(MainColor())

	coverbox := tview.NewTextView()
	coverbox.SetBorder(false).SetBackgroundColor(BackgroundColor())
	coverbox.SetBorderPadding(1, 1, 1, 1)

	roottable, roottablewidth := genMenu("launch", "configure")
	activateTable(roottable)

	launchtable, launchtablewidth := genMenu("node", "wallet", "shell")
	prelightTable(launchtable)

	catstable, catstablewidth := genMenu(app.Cats.GetSortedKeys()...)
	prelightTable(catstable)

	menuflex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(roottable, roottablewidth, 1, true).
		AddItem(coverbox, 0, 1, false)
	menuflex.Box.SetBackgroundColor(BackgroundColor())

	roottable.SetSelectionChangedFunc(func(y, x int) {
		menuflex.
			RemoveItem(coverbox).
			RemoveItem(launchtable).
			RemoveItem(catstable)
		switch y {
		case 0:
			menuflex.
				AddItem(coverbox, 0, 1, true)
		case 1:
			prelightTable(launchtable)
			menuflex.
				AddItem(launchtable, launchtablewidth, 1, true).
				AddItem(coverbox, 0, 1, true)
		case 2:
			menuflex.
				AddItem(catstable, catstablewidth, 1, true).
				AddItem(coverbox, 0, 1, true)
		}
	})
	roottable.SetSelectedFunc(func(y, x int) {
		switch y {
		case 0:
			tapp.Stop()
		case 1:
			activatedTable(roottable)
			activateTable(launchtable)
			tapp.SetFocus(launchtable)
		case 2:
			activatedTable(roottable)
			activateTable(catstable)
			tapp.SetFocus(catstable)
		}
	})

	launchtable.SetSelectionChangedFunc(func(y, x int) {
		switch y {
		case 0:
			coverbox.SetText("")
		case 1:
			coverbox.SetText("run a full peer to peer parallelcoin node")
		case 2:
			coverbox.SetText("\nrun a wallet server (requires a full node)")
		case 3:
			coverbox.SetText("\n\nrun a combined wallet/full node")
		}
	})
	launchtable.SetSelectedFunc(func(y, x int) {
		switch y {
		case 0:
			prelightTable(launchtable)
			activateTable(roottable)
			tapp.SetFocus(roottable)
		}
	})

	var cattable *tview.Table
	var cattablewidth int

	var activepage *tview.Flex

	catstable.SetSelectionChangedFunc(func(y, x int) {
		menuflex.
			RemoveItem(coverbox).
			RemoveItem(cattable)
		if y == 0 {
			menuflex.
				AddItem(coverbox, 0, 1, true)
			return
		}
		cat := app.Cats.GetSortedKeys()[y-1]
		cattable, cattablewidth = genMenu(app.Cats[cat].GetSortedKeys()...)
		prelightTable(cattable)
		cattable.SetSelectedFunc(func(y, x int) {
			if y == 0 {
				activatedTable(roottable)
				prelightTable(cattable)
				activateTable(catstable)
				tapp.SetFocus(catstable)
			}
		})
		cattable.SetSelectionChangedFunc(func(y, x int) {
			menuflex.
				RemoveItem(coverbox).
				RemoveItem(activepage)
			if y == 0 {
				menuflex.AddItem(coverbox, 0, 1, false)
			} else {
				itemname := app.Cats[cat].GetSortedKeys()[y-1]
				activepage =
					genPage(cat, itemname, app)
				menuflex.AddItem(activepage, 0, 1, true)
			}
		})
		menuflex.AddItem(cattable, cattablewidth, 1, false).
			AddItem(coverbox, 0, 1, true)
	})
	catstable.SetSelectedFunc(func(y, x int) {
		if y == 0 {
			prelightTable(catstable)
			activateTable(roottable)
			tapp.SetFocus(roottable)
		} else {
			prelightTable(roottable)
			activatedTable(catstable)
			activateTable(cattable)
			tapp.SetFocus(cattable)
		}
	})

	// root is the canvas (the whole current terminal view)
	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titlebar, 1, 0, false).
		AddItem(menuflex, 0, 1, true)

	if e := tapp.SetRoot(root, true).Run(); e != nil {
		panic(e)
	}

	return 0
}

func genPage(cat, item string, app *config.App) (out *tview.Flex) {
	out = tview.NewFlex().SetDirection(tview.FlexRow)
	out.SetBorderPadding(1, 1, 1, 1)
	out.SetBackgroundColor(PrelightColor())
	// out.SetBorder(true)
	titleblock := tview.NewTextView()
	// titleblock.SetBorder(true)
	titleblock.SetBorderPadding(0, 0, 1, 1)
	titleblock.SetWordWrap(true)
	titleblock.SetBackgroundColor(PrelightColor())
	titleblock.SetTextColor(MainColor())
	titleblock.SetText(
		fmt.Sprintf("%s:%s", strings.ToUpper(cat), strings.ToUpper(item)))
	infoblock := tview.NewTextView()
	// infoblock.SetBorder(true)
	infoblock.SetWordWrap(true)
	infoblock.SetBorderPadding(0, 0, 1, 1)
	infoblock.SetBackgroundColor(PrelightColor())
	infoblock.SetTextColor(MainColor())
	def := app.Cats[cat][item].Default
	defstring := ""
	if def != nil {
		defstring = fmt.Sprintf("\n\ndefault value: %v", def.Get())
	} else {
		defstring = "\n\nthis value has no default"
	}
	infoblock.SetText(
		fmt.Sprintf(
			"%v%s",
			app.Cats[cat][item].Usage, defstring,
		))
	out.AddItem(titleblock, 2, 0, false)
	out.AddItem(infoblock, 5, 0, false)
	row := app.Cats[cat][item]
	switch row.Type {
	case "string", "int", "float", "duration", "port":
		iteminput := tview.NewInputField()
		iteminput.SetBackgroundColor(MainColor())
		iteminput.SetFieldTextColor(PrelightColor())
		iteminput.SetFieldBackgroundColor(MainColor())
		iteminput.SetBorderPadding(1, 1, 1, 1)
		val := app.Cats[cat][item].Value
		if val != nil {
			iteminput.SetText(fmt.Sprint(val.Get()))
		}
		out.AddItem(iteminput, 3, 0, false)
	}
	out.AddItem(tview.NewTextView().SetBackgroundColor(PrelightColor()), 0, 1, false)
	return
}

func getMaxWidth(ss []string) (maxwidth int) {
	for _, x := range ss {
		if len(x) > maxwidth {
			maxwidth = len(x)
		}
	}
	return
}

func genMenu(items ...string) (table *tview.Table, menuwidth int) {
	menuwidth = getMaxWidth(items)
	table = tview.NewTable().SetSelectable(true, true)
	table.SetCell(0, 0, tview.NewTableCell("<"))
	for i, x := range items {
		pad := strings.Repeat(" ", menuwidth-len(x))
		table.SetCell(i+1, 0, tview.NewTableCell(" "+pad+x))
	}
	t, l, _, h := table.Box.GetRect()
	menuwidth += 2
	table.Box.SetRect(t, l, menuwidth, h)
	return
}

// This sets a menu to active attributes
func activateTable(table *tview.Table) {
	rowcount := table.GetRowCount()
	for i := 0; i < rowcount; i++ {
		table.GetCell(i, 0).
			SetAttributes(tcell.AttrNone).
			SetTextColor(TextColor()).
			SetBackgroundColor(MainColor())
		table.SetSelectedStyle(MainColor(), TextColor(), tcell.AttrBold)
		table.Box.SetBackgroundColor(MainColor())
	}
}

// This sets a menu to activated (it has a selected item active)
func activatedTable(table *tview.Table) {
	rowcount := table.GetRowCount()
	for i := 0; i < rowcount; i++ {
		table.GetCell(i, 0).
			SetAttributes(tcell.AttrNone).
			SetTextColor(MainColor()).
			SetBackgroundColor(DimColor())
		table.SetSelectedStyle(DimColor(), MainColor(), tcell.AttrBold)
		table.Box.SetBackgroundColor(DimColor())
	}
}

// This sets a menu to preview (when it is active but not selected yet)
func prelightTable(table *tview.Table) {
	rowcount := table.GetRowCount()
	for i := 0; i < rowcount; i++ {
		table.GetCell(i, 0).
			// SetAttributes(tcell.AttrDim).
			SetTextColor(DimColor()).
			SetBackgroundColor(PrelightColor())
		table.SetSelectedStyle(PrelightColor(), DimColor(), tcell.AttrBold)
		table.Box.SetBackgroundColor(PrelightColor())
	}
}
