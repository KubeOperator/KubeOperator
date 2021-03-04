import {Component, OnInit, ViewChild} from '@angular/core';
import {ThemeService} from "../business/setting/theme/theme.service";
import {HeaderComponent} from "./header/header.component";
import {LicenseService} from "../business/setting/license/license.service";
import {SystemService} from '../business/setting/system.service';

@Component({
    selector: 'app-layout',
    templateUrl: './layout.component.html',
    styleUrls: ['./layout.component.css']
})
export class LayoutComponent implements OnInit {

    constructor(private themeService: ThemeService, private licenseService: LicenseService, private  systemService: SystemService) {
    }

    @ViewChild(HeaderComponent, {static: true})
    header: HeaderComponent;
    alert = false;

    ngOnInit(): void {
        this.licenseService.setLicense();
        this.setTheme();
        this.systemService.getRegistry().subscribe(res => {
            if (res === null) {
                this.alert = true;
            }
            if (res.total === 0) {
                this.alert = true;
            }
        }, error => {
        });
    }

    setTheme() {
        this.themeService.get().subscribe(data => {
            if (data.systemName) {
                document.title = data.systemName;
            }
            if (data.logo) {
                this.header.setLogo(data.logo);
            }
        });
    }
}
