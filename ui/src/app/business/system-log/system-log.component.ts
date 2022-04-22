import {Component, OnInit} from '@angular/core';
import {SystemLogService} from './system-log.service';
import {TranslateService} from '@ngx-translate/core';
import { ClrDatagridStateInterface } from '@clr/angular';

@Component({
    selector: 'app-system-log',
    templateUrl: './system-log.component.html',
    styleUrls: ['./system-log.component.css']
})
export class SystemLogComponent implements OnInit {
    loading = false;
    total = 0;
    page = 1;
    size = 15;
    items = [];
    defaultFilter = 'name';
    currentTerm: string = '';
    isOpenFilterTag: boolean;
    constructor(private service: SystemLogService, public translate: TranslateService) {}

    ngOnInit(): void {
        document.oncontextmenu =function () {return false; };
        document.onkeydown=function(){
            var e=window.event||arguments[0];
            if(e.ctrlKey||e.keyCode==83){
                return false;
            }
        };
    }
    
    get inProgress(): boolean {
        return this.loading;
    }

    public doFilter(terms: string): void {
        // allow search by null characters
        if (terms === undefined || terms === null) {
            return;
        }
        this.currentTerm = terms.trim();
        this.loading = true;
        this.page = 1;
        this.total = 0;
        this.load();
    }

    refresh(): void {
        this.doFilter("");
    }

    filter() {
        this.load();
    }

    openFilter(isOpen: boolean) {
        this.isOpenFilterTag = isOpen;
    }

    selectFilterKey($event: any): void {
        this.defaultFilter = $event['target'].value;
        this.doFilter(this.currentTerm);
    }

    load(state?: ClrDatagridStateInterface) {
        if (state && state.page) {
            this.size = state.page.size;
        }
        this.loading = true;
        this.service.list(this.page, this.size, this.defaultFilter, this.currentTerm).subscribe(data => {
            const currentLanguage = localStorage.getItem('currentLanguage') || this.translate.getBrowserCultureLang();
            this.items = data.items;
            if (this.items != null) {
                for (const item of this.items) {
                    item.isExceed = false;
                    if (currentLanguage == 'en-US') {
                        item.operation = item.operation.split('|')[1];
                    } else {
                        item.operation = item.operation.split('|')[0];
                    }
                    if (item.operation.length > 40) {
                        item.isExceed = true;
                        item.completeOperation = item.operation;
                        item.operation = item.operation.substring(0,40);
                    }
                }
            }
            this.total = data.total;
            this.loading = false;
        });
        this.loading = false;
    }
}
