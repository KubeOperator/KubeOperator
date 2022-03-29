import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Cluster} from '../../../cluster';
import {MonitorService} from '../monitor.service';
import {DomSanitizer} from '@angular/platform-browser';
import {ClusterTool} from "../../tools/tools";
import * as echarts from 'echarts';
import { NodeService } from '../../node/node.service';

@Component({
    selector: 'app-monitor-dashboard',
    templateUrl: './monitor-dashboard.component.html',
    styleUrls: ['./monitor-dashboard.component.css']
})
export class MonitorDashboardComponent implements OnInit {

    loading: boolean = false;
    
    @Input() currentCluster: Cluster;
    @Input() item: ClusterTool;
    nodes: string[] = [];
    selectNode: string = '';
    searchBeginDate: Date;
    searchEndDate: Date;
    cpuDateList: string[];
    cpuValueList: any[] = [];
    memeryDateList: string[];
    memeryValueList: any[] = [];
    diskDateList: string[];
    diskValueList: any[] = [];
    networkDateList: string[];
    networkValueList: any[] = [];

    constructor(private nodeService: NodeService, private service: MonitorService, private sanitizer: DomSanitizer) {
    }

    ngOnInit(): void {
        this.nodeService.listWithoutPage(this.currentCluster.name).subscribe(data => {
            this.nodes = data.items.map(function (item) {
                return item.ip;
            })
            this.selectNode = this.nodes[0]
            this.load();
        }, error => {
            this.loading = false;
            // this.alertService.showAlert(error.error.msg, AlertLevels.ERROR);
        })
    }

    load() {
        const begainDate = (this.searchBeginDate === undefined) ? new Date(new Date().setMinutes(new Date().getMinutes() - 30)) : new Date(this.searchBeginDate);
        const endDate = (this.searchEndDate === undefined) ? new Date() : new Date(this.searchEndDate);

        const start = begainDate.getTime() / 1000;
        const end = endDate.getTime() / 1000;
        
        this.getCPUDatas(start, end);
        this.getMemeryDatas(start, end);
        this.getDiskDatas(start, end);
        this.getNetworkDatas(start, end);
    }


    getCPUDatas (start: Number, end: Number) {
        this.cpuDateList = [];
        this.cpuValueList = [];
        this.service.QueryCPU(this.currentCluster.name, (this.selectNode + ':9100'), '"system"', start.toString(), end.toString()).subscribe(data => {
            if (data.data.result.length === 0) {
                this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
                return;
            }
            this.cpuDateList = data.data.result[0].values.map(function (item) {
                const timeNow = new Date(item[0] * 1000);
                return (timeNow.getMonth() + 1) + "月" + timeNow.getDate() + "日" + timeNow.getHours() + ":" + timeNow.getMinutes();
            })
            let itemDatas: string[]
            itemDatas = data.data.result[0].values.map(function (item) {
                return Number(item[1]).toFixed(2);
            })
            this.cpuValueList.push(this.addSeries(itemDatas, 'Busy System'));
            this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
            
            this.service.QueryCPU(this.currentCluster.name, (this.selectNode + ':9100'), '"user"', start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[];
                itemDatas = data.data.result[0].values.map(function (item) {
                    return Number(item[1]).toFixed(2);
                })
                this.cpuValueList.push(this.addSeries(itemDatas, 'Busy User'));
                this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
            })

            this.service.QueryCPU(this.currentCluster.name, (this.selectNode + ':9100'), '"iowait"', start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[];
                itemDatas = data.data.result[0].values.map(function (item) {
                    return Number(item[1]).toFixed(2);
                })
                this.cpuValueList.push(this.addSeries(itemDatas, 'Busy Iowait'));
                this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
            })

            this.service.QueryCPU(this.currentCluster.name, (this.selectNode + ':9100'), '"idle"', start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[];
                itemDatas = data.data.result[0].values.map(function (item) {
                    return Number(item[1]).toFixed(2);
                })
                this.cpuValueList.push(this.addSeries(itemDatas, 'Busy Idle'));
                this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
            })

            this.service.QueryCPU(this.currentCluster.name, (this.selectNode + ':9100'), '~".*irq"', start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[];
                itemDatas = data.data.result[0].values.map(function (item) {
                    return Number(item[1]).toFixed(2);
                })
                this.cpuValueList.push(this.addSeries(itemDatas, 'Busy Irqs'));
                this.initCharts('cpuChart', 'CPU Basic', this.cpuDateList, this.cpuValueList, '%');
            })
        })
    }

    getMemeryDatas (start: Number, end: Number) {
        this.memeryDateList = [];
        this.memeryValueList = [];
        this.service.QueryMemeryTotal(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
            if (data.data.result.length === 0) {
                this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
                return;
            }
            this.memeryDateList = data.data.result[0].values.map(function (item) {
                const timeNow = new Date(item[0] * 1000);
                return (timeNow.getMonth() + 1) + "月" + timeNow.getDate() + "日" + timeNow.getHours() + ":" + timeNow.getMinutes();
            })
            let itemDatas: string[]
            itemDatas = data.data.result[0].values.map(function (item) {
                return (Number(item[1]) / 1024 / 1024 / 1024).toFixed(2);
            })

            this.service.QueryMemeryUsed(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[]
                itemDatas = data.data.result[0].values.map(function (item) {
                    return (Number(item[1]) / 1024 / 1024 / 1024).toFixed(2);
                })
                this.memeryValueList.push(this.addSeries(itemDatas, 'RAM Used'));
                this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
            })
            this.service.QueryMemeryCacheBuffer(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[]
                itemDatas = data.data.result[0].values.map(function (item) {
                    return (Number(item[1]) / 1024 / 1024 / 1024).toFixed(2);
                })
                this.memeryValueList.push(this.addSeries(itemDatas, 'RAM Cache + Buffer'));
                this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
            })
            this.service.QueryMemeryFree(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[]
                itemDatas = data.data.result[0].values.map(function (item) {
                    return (Number(item[1]) / 1024 / 1024 / 1024).toFixed(2);
                })
                this.memeryValueList.push(this.addSeries(itemDatas, 'RAM Free'));
                this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
            })
            this.service.QueryMemerySWAPUsed(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
                let itemDatas: string[]
                itemDatas = data.data.result[0].values.map(function (item) {
                    return (Number(item[1]) / 1024 / 1024 / 1024).toFixed(2);
                })
                this.memeryValueList.push(this.addSeries(itemDatas, 'SWAP Used'));
                this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
            })

            this.memeryValueList.push(this.addSeries(itemDatas, 'RAM Total'));
            this.initCharts('memeryChart', 'Memery Basic', this.memeryDateList, this.memeryValueList, 'GiB');
        })
    }

    getDiskDatas (start: Number, end: Number) {
        this.diskDateList = [];
        this.diskValueList = [];
        this.service.QueryDisk(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
            if (data.data.result.length === 0) {
                this.initCharts('diskChart', 'Disk Space Used Basic', this.diskDateList, this.diskValueList, '%');
                return;
            }
            this.diskDateList = data.data.result[0].values.map(function (item) {
                const timeNow = new Date(item[0] * 1000);
                return (timeNow.getMonth() + 1) + "月" + timeNow.getDate() + "日" + timeNow.getHours() + ":" + timeNow.getMinutes();
            })
            let itemDatas: string[];
            itemDatas = data.data.result[0].values.map(function (item) {
                return Number(item[1]).toFixed(2);
            })
            this.diskValueList.push(this.addSeries(itemDatas, 'Disk Space Used'));
            this.initCharts('diskChart', 'Disk Space Used Basic', this.diskDateList, this.diskValueList, '%');
        })
    }

    getNetworkDatas (start: Number, end: Number) {
        this.networkDateList = [];
        this.networkValueList = [];
        this.service.QueryNetworkRecv(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
            if (data.data.result.length === 0) {
                this.initCharts('networkChart', 'Network Traffic Basic', this.networkDateList, this.networkValueList, 'kb/s');
                return;
            }
            this.networkDateList = data.data.result[0].values.map(function (item) {
                const timeNow = new Date(item[0] * 1000);
                return (timeNow.getMonth() + 1) + "月" + timeNow.getDate() + "日" + timeNow.getHours() + ":" + timeNow.getMinutes();
            })
            for (const res of data.data.result) {
                let itemDatas: string[];
                itemDatas = res.values.map(function (item) {
                    return (Number(item[1]) / 1000).toFixed(0);
                })
                this.networkValueList.push(this.addSeries(itemDatas, 'Recv ' + res.metric.device));
                this.initCharts('networkChart', 'Network Traffic Basic', this.networkDateList, this.networkValueList, 'kb/s');
            }

            this.service.QueryNetworkTrans(this.currentCluster.name, (this.selectNode + ':9100'), start.toString(), end.toString()).subscribe(data => {
                for (const res of data.data.result) {
                    let itemDatas: string[];
                    itemDatas = res.values.map(function (item) {
                        return -(Number(item[1]) / 1000).toFixed(0);
                    })
                    this.networkValueList.push(this.addSeries(itemDatas, 'Trans ' + res.metric.device));
                    this.initCharts('networkChart', 'Network Traffic Basic', this.networkDateList, this.networkValueList, 'kb/s');
                }
            })
        })
    }

    initCharts(chartName: string, title: string, xDatas: string[], yDatas: string[], formatStr: string) {
        const lineChart = echarts.init(document.getElementById(chartName));
        const option = {
            title: [{
                left: 'center',
                text: title
            }],
            tooltip: {
                trigger: 'axis',
                formatter: function (datas) {
                    var res = datas[0].name + '<br/>';
                    for (const item of datas) {
                       res += item.marker + " " + item.seriesName + '：' + item.data + formatStr + '<br/>';
                    }
                    return res;
                },
                textStyle:{
                    align:'left'
                }
            },
            xAxis: [{
                data: xDatas,
                gridIndex: 1
            }],
            yAxis: [{
                axisLabel: {formatter: '{value} '+ formatStr},
                gridIndex: 1
            }],
            grid: [{bottom: '60%'}, {left: '20%'}, {top: '15%'}],
            series: yDatas
        };
        lineChart.setOption(option, true);
    }

    addSeries (datas: string[], name: string):any {
        return {
            name: name,
            type: 'line',
            smooth: true,
            showSymbol: true,
            areaStyle: {},
            data: datas
        }
    }
}
