import {Component, OnInit, Input} from '@angular/core';
import {LoggingService} from '../logging.service';
import {LokiData} from './loki';

@Component({
    selector: 'app-loki-logging',
    templateUrl: './loki.component.html',
    styleUrls: ['./loki.component.css']
})
export class LokiComponent implements OnInit {

    @Input() clusterName: string;

    constructor(private service: LoggingService) {
    }

    loading = false;
    searchBeginDate: Date;
    searchEndDate: Date;
    labels: string[] = [];
    label: string = '';
    values: string[] = [];
    value: string = '';
    
    logs: LokiData[]= [];

    ngOnInit(): void {
        this.getLabels();
    }
    getLabels () {
        this.service.LokiLabels(this.clusterName).subscribe(data => {
            this.labels = data.data;
        });
    }
    getLabelValues () {
        this.service.LokiLabelValues(this.clusterName, this.label).subscribe(data => {
            if (data.data === undefined) {
                this.values = [];
                this.value = '';
            } else {
                if (data.data.length > 0) {
                    this.values = data.data;
                    this.value = this.values[0];
                } else {
                    this.values = [];
                    this.value = '';
                }
            }
        });
    }
    
    refresh() {
        this.loading = true;
        this.logs = [];
        let step: number = 8;
        let paramInfo = 'direction=backward&limit=1000&regexp=';
        let start: number = new Date(this.searchBeginDate).getTime();
        let end: number = new Date(this.searchEndDate).getTime();
        if (!isNaN(start) && !isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + start + '000000&end=' + (end + 86400000) + '000000');
            step = (((end - start) / 86400000) + 1) * 8;
        }
        if (isNaN(start) && !isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + end + '000000&end=' + (end + 86400000) + '000000');
        }
        if (!isNaN(start) && isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + start + '000000&end=' + (start + 86400000) + '000000');
        }
        paramInfo = paramInfo + '&step=' + step;

        if (this.label !== '' && this.value !== '') {
            paramInfo += ('&query={' + this.label + '="' + this.value + '"}');
        }
        this.service.LokiSearch(this.clusterName, paramInfo).subscribe(data => {
            for (const item1 of data.data.result) {
                for (const item2 of item1.values) {
                    let logItem = this.dataParser(item2);
                    this.logs.push(logItem);
                }
            }
            this.loading = false;
        }, error => {
            this.loading = false;
        });
    }
    dataParser (data: any): LokiData{
        let logItem:LokiData= {
            ts: new Date(parseInt(data[0].substring(0 ,13))),
            info: data[1],
        }
        return logItem;
    }
}
