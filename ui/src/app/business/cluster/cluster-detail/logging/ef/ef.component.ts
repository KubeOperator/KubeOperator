import {Component, OnInit, Input} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {LoggingService} from '../logging.service';
import format from 'date-fns/format';
import {ToolsService} from '../../tools/tools.service';

@Component({
    selector: 'app-ef-logging',
    templateUrl: './ef.component.html',
    styleUrls: ['./ef.component.css']
})
export class EfComponent implements OnInit {

    @Input() clusterName: string;
    @Input() indexPrefix: string;

    constructor(private service: LoggingService, private toolsService: ToolsService, private route: ActivatedRoute) {
    }

    loading = false;
    total = 0;
    page = 1;
    size = 10;
    isExpertShow = false;
    searchInfo = '';
    searchBeginDate: Date;
    searchEndDate: Date;
    namespace = '';
    pod = '';
    container = '';
    logs = [];

    ngOnInit(): void {
        this.refresh();
    }
    changeSearchType(searchType) {
        this.isExpertShow = !searchType;
    }
    refresh() {
        const queryArry = [];
        let beginDate;
        let endDate;
        let queryIndex = '';
        if (this.searchInfo !== '') {
            queryArry.push({match: {message: {query: this.searchInfo}}});
        }
        if (this.searchBeginDate === undefined) {
            beginDate = format(new Date(), 'yyyy.MM.dd');
        } else {
            beginDate = format(new Date(this.searchBeginDate), 'yyyy.MM.dd');
        }
        if (this.searchEndDate === undefined) {
            endDate = format(new Date(), 'yyyy.MM.dd');
        } else {
            endDate = format(new Date(this.searchEndDate), 'yyyy.MM.dd');
        }
        const itemDate = new Date(beginDate);
        while (itemDate <= new Date(endDate)){
            queryIndex = queryIndex + this.indexPrefix + format(itemDate, 'yyyy.MM.dd') + ',';
            itemDate.setDate(itemDate.getDate() + 1);
        }
        queryIndex = queryIndex.substr(0, queryIndex.length - 1);
        if (this.isExpertShow) {
            if (this.namespace !== '') {
                if (this.namespace.indexOf('-') === -1) {
                    queryArry.push({term: {'kubernetes.namespace_name': this.namespace}});
                } else {
                    const namespaceArry = this.namespace.split('-');
                    for (const n of namespaceArry) {
                        queryArry.push({term: {'kubernetes.namespace_name': n}});
                    }
                }
            }
            if (this.pod !== '') {
                if (this.pod.indexOf('-') === -1) {
                    queryArry.push({term: {'kubernetes.pod_name': this.pod}});
                } else {
                    const podArry = this.pod.split('-');
                    for (const p of podArry) {
                        queryArry.push({term: {'kubernetes.pod_name': p}});
                    }
                }
            }
            if (this.container !== '') {
                if (this.container.indexOf('-') === -1) {
                    queryArry.push({term: {'kubernetes.container_name': this.container}});
                } else {
                    const containerArry = this.container.split('-');
                    for (const c of containerArry) {
                        queryArry.push({term: {'kubernetes.container_name': c}});
                    }
                }
            }
        }
        this.service.EfSearch(this.clusterName, queryArry, queryIndex, beginDate, endDate, this.page, this.size).subscribe(data => {
            this.logs = data.hits.hits;
            this.total = data.hits.total.value;
            for (const item of this.logs) {
                const timeItem = new Date(item._source['@timestamp']);
                item.timestamp = timeItem.getFullYear() + '-' + (timeItem.getMonth() + 1) + '-' + timeItem.getDate() + ' ' +
                    timeItem.getHours() + ':' + timeItem.getMinutes() + ':' + timeItem.getSeconds();
                item._source = JSON.stringify(item._source);
            }
            this.loading = false;
        });
    }
}
