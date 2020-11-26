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
            this.labels = data.data
        });
    }
    getLabelValues () {
        this.service.LokiLabelValues(this.clusterName, this.label).subscribe(data => {
            this.values = data.data
        });
    }
    
    refresh() {
        this.loading = true;
        let paramInfo = 'limit=1000';
        let start: number = new Date(this.searchBeginDate).getTime();
        let end: number = new Date(this.searchEndDate).getTime();
        if (!isNaN(start) && !isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + start + '&end=' + (end + 86400000))
        }
        if (isNaN(start) && !isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + end + '&end=' + end + 86400000)
        }
        if (!isNaN(start) && isNaN(end)) {
            paramInfo = paramInfo + ('&start=' + start + '&end=' + (start + 86400000))
        }

        if (this.label !== '' && this.value !== '') {
            paramInfo += ('&query={' + this.label + '="' + this.value + '"}')
        }
        this.service.LokiSearch(this.clusterName, paramInfo).subscribe(data => {
            for (let i = 0; i < data.data.result.length; i++) {
                for (let j = 0; j < data.data.result[i].values.length; j++) {
                    let logItem = this.dataParser(data.data.result[i].values[j][1])
                    this.logs.push(logItem)
                }
            }
            this.loading = false;
        }, error => {
            this.loading = false;
        });
    }
    dataParser (Str): LokiData{
        let logItem:LokiData= {
            ts: new Date,
            info: Str,
        }
        let dataArry = Str.split(' ')
        for (let i = 0; i < dataArry.length; i++) {
            if (dataArry[i].indexOf('ts=') !== -1) {
                logItem.ts = new Date(dataArry[i].split('=')[1])
                break
            } else if (dataArry[i].indexOf('[') !== -1 && dataArry[i].indexOf(']') !== -1) {
                logItem.ts = new Date(dataArry[i].replace('[', '').replace(']', ''))
                break
            }
        }
        return logItem
    }
}
