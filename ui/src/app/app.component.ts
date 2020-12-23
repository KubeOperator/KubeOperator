import {Component, OnInit} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';


@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

    title = 'app';

    constructor(public translate: TranslateService) {
    }

    ngOnInit() {
        const currentLanguage = localStorage.getItem('currentLanguage') || this.translate.getBrowserCultureLang();
        this.translate.setDefaultLang('zh-CN');
        this.translate.use(currentLanguage);
        localStorage.setItem('currentLanguage', currentLanguage);
    }


}
