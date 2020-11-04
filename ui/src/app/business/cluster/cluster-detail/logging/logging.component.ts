import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {LoggingService} from './logging.service';
import format from 'date-fns/format';
import {ToolsService} from '../tools/tools.service';

@Component({
    selector: 'app-logging',
    templateUrl: './logging.component.html',
    styleUrls: ['./logging.component.css']
})
export class LoggingComponent implements OnInit {

    constructor(private service: LoggingService, private toolsService: ToolsService, private route: ActivatedRoute) {
    }

    loading = false;
    total = 0;
    page = 1;
    size = 10;
    isExpertShow = false;
    searchIndexPrefix = '';
    searchInfo = '';
    searchBeginDate: Date;
    searchEndDate: Date;
    namespace = '';
    pod = '';
    container = '';
    currentCluster: Cluster;
    logs = [];

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolsService.list(this.currentCluster.name).subscribe(tools => {
                for (const tool of tools) {
                    if (tool.name === 'logging') {
                        this.searchIndexPrefix = tool.vars['fluentd-elasticsearch.elasticsearch.logstashPrefix'] + '-';
                        this.loading = true;
                        this.refresh();
                        break;
                    }
                }
            });
        });
    }
    changeSearchType(searchType) {
        this.isExpertShow = !searchType;
    }
    refresh() {
        const queryArry = [];
        let beginDate = '';
        let endDate = '';
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
            queryIndex = queryIndex + this.searchIndexPrefix + format(itemDate, 'yyyy.MM.dd') + ',';
            itemDate.setDate(itemDate.getDate() + 1);
        }
        queryIndex = queryIndex.substr(0, queryIndex.length - 1);
        if (this.isExpertShow) {
            if (this.namespace !== '') {
                queryArry.push({match: {'kubernetes.namespace_name': {query: this.namespace}}});
            }
            if (this.pod !== '') {
                queryArry.push({match: {'kubernetes.pod_name': {query: this.pod}}});
            }
            if (this.container !== '') {
                queryArry.push({match: {'kubernetes.container_name': {query: this.container}}});
            }
        }
        this.service.Search(this.currentCluster.name, queryArry, queryIndex, beginDate, endDate, this.page, this.size).subscribe(data => {
            this.logs = data.hits.hits;
            this.total = data.hits.total.value;
            for (const item of this.logs) {
                const timeItem = new Date(item._source['@timestamp']);
                item.timestamp = timeItem.getFullYear() + '-' + (timeItem.getMonth() + 1) + '-' + timeItem.getDate() + ' ' +
                    timeItem.getHours() + ':' + timeItem.getMinutes() + ':' + timeItem.getSeconds();
                const sourceStr = JSON.stringify(item._source);
                item._source = sourceStr;
            }
            this.loading = false;
        });
    }
}
