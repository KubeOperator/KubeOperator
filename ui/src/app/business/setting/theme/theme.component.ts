import {Component, OnInit} from '@angular/core';
import {Theme} from "./theme";
import {ThemeService} from "./theme.service";
import {CommonAlertService} from "../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../layout/common-alert/alert";
import {TranslateService} from "@ngx-translate/core";

@Component({
    selector: 'app-theme',
    templateUrl: './theme.component.html',
    styleUrls: ['./theme.component.css']
})
export class ThemeComponent implements OnInit {

    constructor(private themeService: ThemeService, private alertService: CommonAlertService, private translateService: TranslateService) {
    }

    item: Theme = new Theme();
    file: File;
    logoBase64: string;

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.themeService.get().subscribe(data => {
            this.item = data;
        });
    }

    onLogoChange(e: any) {
        this.file = e.target.files[0];
        const r = new FileReader();
        r.readAsDataURL(e.target.files[0]);
        r.onload = (b) => {
            this.logoBase64 = b.target.result.toString();
        };
    }

    onSubmit() {
        this.item.logo = this.logoBase64;
        this.themeService.update(this.item).subscribe(data => {
            this.refresh();
            this.alertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            window.location.reload();
        }, error => {
            this.alertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onCancel() {
        this.refresh();
    }

}
