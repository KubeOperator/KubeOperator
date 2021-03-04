import {Component, OnInit, ViewChild} from '@angular/core';
import {ThemeService} from '../business/setting/theme/theme.service';
import {HeaderComponent} from './header/header.component';
import {LicenseService} from '../business/setting/license/license.service';
import {RegistryService} from '../business/setting/registry-setting/registry.service';

@Component({
    selector: 'app-layout',
    templateUrl: './layout.component.html',
    styleUrls: ['./layout.component.css']
})
export class LayoutComponent implements OnInit {

    constructor(private themeService: ThemeService, private licenseService: LicenseService, private  registryService: RegistryService) {
    }

    @ViewChild(HeaderComponent, {static: true})
    header: HeaderComponent;
    alert = false;

    ngOnInit(): void {
        this.licenseService.setLicense();
        this.setTheme();
        this.registryService.mixedGet(1, 1).subscribe(registry => {
            if (registry.total < 1) {
                this.alert = true;
            }
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
