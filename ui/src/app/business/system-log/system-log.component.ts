import {Component, OnInit} from '@angular/core';
import {SystemLogService} from './system-log.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-system-log',
    templateUrl: './system-log.component.html',
    styleUrls: ['./system-log.component.css']
})
export class SystemLogComponent implements OnInit {
    loading = false;
    total = 0;
    page = 1;
    size = 10;
    items = [];
    constructor(private service: SystemLogService, public translate: TranslateService) {}

    ngOnInit(): void {
        this.refresh()
    }
    refresh() {
        this.loading = true;
        this.service.list(this.page, this.size).subscribe(data => {
            const currentLanguage = localStorage.getItem('currentLanguage') || this.translate.getBrowserCultureLang();
            this.items = data.items;
            if (this.items != null) {
                for (const item of this.items) {
                    if (currentLanguage == 'en-US') {
                        item.operation = item.operation.split('|')[1]
                        item.operationUnit = item.operationUnit.split('|')[1]
                    } else {
                        item.operation = item.operation.split('|')[0]
                        item.operationUnit = item.operationUnit.split('|')[0]
                    }
                }
            }
            this.total = data.total;
            this.loading = false;
        });
    }
}
