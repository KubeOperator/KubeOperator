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
        this.setTheme();
    }

    setTheme() {
        this.licenseService.get().subscribe(d => {
            sessionStorage.setItem('license', JSON.stringify(d));
            this.themeService.get().subscribe(data => {
                if (data.systemName) {
                    document.title = data.systemName;
                }
                if (data.logo) {
                    this.header.setLogo(data.logo);
                }
            });
        });

    }
}
