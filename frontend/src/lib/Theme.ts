
/*

From Ghost Branding Guidelines

charcoal:      #222531
midnight blue: #492CFB
lapis:         #4A64FA
red:           #ED475B
green:         #93D65E
orange:        #F5A42A
blue:          #58BAC8

steel grey:    #6B7292
cloud grey:    #C3C8DF
light grey:    #F2F4F7

 */

import {ThemeDefinition} from "vuetify";

const ghostTheme: ThemeDefinition = {
    dark: true,
    colors: {
        background: '#333747',
        surface: '#222531',
        'background-lighten-1': '#6B7292',
        primary: '#492CFB',
        'primary-lighten-1': '#4A64FA',
        secondary: '#FA7400',
        'secondary-darken-1': '#FF8922',
        error: '#ED475B',
        info: '#58BAC8',
        success: '#93D65E',
        warning: '#F5A42A',
    }
}

export default ghostTheme;