import {Component, OnInit, ViewChild} from '@angular/core';
import {ThemeService} from "../business/setting/theme/theme.service";
import {Theme} from "../business/setting/theme/theme";
import {HeaderComponent} from "./header/header.component";
import {LicenseService} from "../business/setting/license/license.service";

@Component({
    selector: 'app-layout',
    templateUrl: './layout.component.html',
    styleUrls: ['./layout.component.css']
})
export class LayoutComponent implements OnInit {

    constructor(private themeService: ThemeService, private licenseService: LicenseService) {
    }

    @ViewChild(HeaderComponent, {static: true})
    header: HeaderComponent;

    ngOnInit(): void {
        this.licenseService.setLicense();

        this.themeService.setTheme();
        const str = sessionStorage.getItem('theme');
        const theme: Theme = JSON.parse(str);
        if (theme.systemName) {
            document.title = theme.systemName;
        }
        if (theme.logo) {
            this.header.setLogo(theme.logo);
        }
    }
}
