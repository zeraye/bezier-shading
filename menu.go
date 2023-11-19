package bezier_shading

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zeraye/bezier-shading/pkg/config"
	"github.com/zeraye/bezier-shading/pkg/draw"
)

type Menu struct {
	config                     *config.Config
	kdSlider                   *widget.Slider
	ksSlider                   *widget.Slider
	mSlider                    *widget.Slider
	lightAnimationButton       *widget.Button
	lightHeightSlider          *widget.Slider
	backgroundSolidColorLabel  *widget.Label
	backgroundSolidColorButton *widget.Button
	backgroundImageLabel       *widget.Label
	backgroundImageButton      *widget.Button
	pointsHeightSlider         *widget.Slider
	pointsHeightContainer      *fyne.Container
}

func NewMenu(config *config.Config) *Menu {
	return &Menu{config: config}
}

func (m *Menu) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(m.config.Window.Width-m.config.UI.RasterWidth), float32(m.config.UI.RasterHeight))

}

func (m *Menu) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// layout for settingsLabel
	topLeft := fyne.NewPos(0, 0)
	objects[0].Resize(size)
	objects[0].Move(topLeft)

	// layout for other objcets
	padding := theme.Padding()
	for _, child := range objects[1:] {
		childMin := child.MinSize()
		childMin.Width = size.Width - 6*padding // magic number, make UI look nice
		child.Resize(childMin)
		child.Move(fyne.NewPos(float32(size.Width-childMin.Width)/2, float32(size.Height-childMin.Height)/2))
	}
}

func (m *Menu) BuildUI(g *Game) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	kdValue := m.config.Defaults.Kd
	kdBinding := binding.BindFloat(&kdValue)
	kdLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(kdBinding, "k_d (%0.2f)"))
	kdSlider := widget.NewSliderWithData(0, 1, kdBinding)
	kdSlider.Step = 0.01
	m.kdSlider = kdSlider

	ksValue := m.config.Defaults.Ks
	ksBinding := binding.BindFloat(&ksValue)
	ksLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(ksBinding, "k_s (%0.2f)"))
	ksSlider := widget.NewSliderWithData(0, 1, ksBinding)
	ksSlider.Step = 0.01
	m.ksSlider = ksSlider

	mValue := m.config.Defaults.M
	mBinding := binding.BindFloat(&mValue)
	mLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(mBinding, "m (%0.0f)"))
	mSlider := widget.NewSliderWithData(1, 100, mBinding)
	mSlider.Step = 1
	m.mSlider = mSlider

	defaultLightColor := draw.RGBAToColor(m.config.Defaults.LightColorRGBA)
	lcr, lcg, lcb, _ := draw.ColorRGBA(defaultLightColor)
	lightColorLabel := widget.NewLabel(fmt.Sprintf("color: (%d, %d, %d)", lcr, lcg, lcb))
	lightColorButton := widget.NewButton("Pick light color", lightColorButtonTapped(g, lightColorLabel))

	lightAnimationButton := widget.NewButton("", animationButtonTapped(g))
	if g.LightAnimation {
		lightAnimationButton.Text = "Stop animation"
	} else {
		lightAnimationButton.Text = "Start animation"
	}
	m.lightAnimationButton = lightAnimationButton

	lightHeightBinding := binding.BindFloat(&g.lightHeight)
	lightHeightLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(lightHeightBinding, "light height (%0.0f)"))
	lightHeightSlider := widget.NewSliderWithData(1, 400, lightHeightBinding)
	lightHeightSlider.Step = 1
	m.lightHeightSlider = lightHeightSlider

	backgroundRadioButton := widget.NewRadioGroup([]string{"Solid color", "Image"}, nil)
	backgroundRadioButton.SetSelected("Solid color")
	backgroundRadioButton.OnChanged = backgroundRadioButtonChanged(g, backgroundRadioButton)

	defaultBackgroundSolidColor := draw.RGBAToColor(m.config.Defaults.DefaultBackgroundSolidColorRGBA)
	bscr, bscg, bscb, _ := draw.ColorRGBA(defaultBackgroundSolidColor)
	backgroundSolidColorLabel := widget.NewLabel(fmt.Sprintf("color: (%d, %d, %d)", bscr, bscg, bscb))
	backgroundSolidColorButton := widget.NewButton("Pick background solid color", backgroundSolidColorButtonTapped(g, backgroundSolidColorLabel))
	m.backgroundSolidColorLabel = backgroundSolidColorLabel
	m.backgroundSolidColorButton = backgroundSolidColorButton

	backgroundImageLabel := widget.NewLabel("file: -")
	backgroundImageButton := widget.NewButton("Open background image file", backgroundImageButtonTapped(g, backgroundImageLabel))
	backgroundImageLabel.Hide()
	backgroundImageButton.Hide()
	m.backgroundImageLabel = backgroundImageLabel
	m.backgroundImageButton = backgroundImageButton

	normalMapLabel := widget.NewLabel("file: -")
	normalMapButton := widget.NewButton("Open normal map file", normalMapButtonTapped(g, normalMapLabel))

	triangulationLabel := widget.NewLabel("triangulation")
	triangulationSlider := widget.NewSlider(2, 29)
	triangulationSlider.Step = 1
	triangulationSlider.OnChanged = triangulationSliderChanged(g, triangulationSlider)
	triangulationSlider.Value = float64(g.config.Defaults.Triangulation)
	triangulationCheck := widget.NewCheck("show mesh", triangulationCheckChanged(g))

	pointsHeightLabel := widget.NewLabel("point height")
	pointsHeightSlider := widget.NewSlider(0, 200)
	pointsHeightSlider.OnChanged = pointsHeightSliderChanged(g, pointsHeightSlider)
	m.pointsHeightSlider = pointsHeightSlider
	pointsHeightContainer := container.NewGridWithColumns(2, pointsHeightLabel, pointsHeightSlider)
	m.pointsHeightContainer = pointsHeightContainer

	return container.New(m, title, container.NewVBox(
		container.NewGridWithColumns(2, kdLabel, kdSlider),
		container.NewGridWithColumns(2, ksLabel, ksSlider),
		container.NewGridWithColumns(2, mLabel, mSlider),
		container.NewGridWithColumns(2, lightColorLabel, lightColorButton),
		container.NewGridWithColumns(2, lightHeightLabel, lightHeightSlider),
		backgroundRadioButton,
		backgroundSolidColorLabel,
		backgroundSolidColorButton,
		backgroundImageLabel,
		backgroundImageButton,
		normalMapLabel,
		normalMapButton,
		container.NewGridWithColumns(3, triangulationLabel, triangulationSlider, triangulationCheck),
		lightAnimationButton,
		pointsHeightContainer,
	))
}
