module github.com/SCKelemen/viz-cli

go 1.25.4

replace github.com/SCKelemen/dataviz => ../dataviz

replace github.com/SCKelemen/design-system => ../design-system

replace github.com/SCKelemen/wpt-test-gen => ../wpt-test-gen

replace github.com/SCKelemen/layout => ../layout

replace github.com/SCKelemen/render-svg => ../render-svg

replace github.com/SCKelemen/color => ../color

replace github.com/SCKelemen/text => ../text

replace github.com/SCKelemen/unicode => ../unicode

replace github.com/SCKelemen/units => ../units

require (
	github.com/SCKelemen/dataviz v0.0.0-00010101000000-000000000000
	github.com/SCKelemen/design-system v0.0.0-20260108142421-70048f811d38
)

require (
	github.com/SCKelemen/color v1.0.0 // indirect
	github.com/SCKelemen/layout v1.1.0 // indirect
	github.com/SCKelemen/render-svg v0.0.0-20260108140101-f0c550c69472 // indirect
	github.com/SCKelemen/text v0.0.0-00010101000000-000000000000 // indirect
	github.com/SCKelemen/unicode v1.0.1-0.20251225190048-233be2b0d647 // indirect
	github.com/SCKelemen/units v1.0.2 // indirect
)
