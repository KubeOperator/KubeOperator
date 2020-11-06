import {Component, OnInit} from '@angular/core';
import {ClusterGradeService} from './cluster-grade.service';
import {Grade, NamespaceResult} from './grade';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-cluster-grade',
    templateUrl: './cluster-grade.component.html',
    styleUrls: ['./cluster-grade.component.css']
})
export class ClusterGradeComponent implements OnInit {

    item: Grade = new Grade();
    currentCluster: Cluster;
    pieChartOptions = {};
    barChartOptions;
    loading = false;

    constructor(private clusterGradeService: ClusterGradeService,
                private route: ActivatedRoute,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.loading = true;
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.clusterGradeService.getGrade(this.currentCluster.name).subscribe(res => {
                this.item = res;
                this.initData(this.item);
                this.item.grade = this.getGrade(this.item.score);
                this.loading = false;
            }, error => {
                this.loading = false;
            });
        });
    }

    initData(item: Grade) {
        this.pieChartOptions = {
            tooltip: {
                trigger: 'item',
                formatter: '{a} <br/>{b}: {c} ({d}%)'
            },
            legend: {
                data: [this.translateService.instant('APP_DANGER'),
                    this.translateService.instant('APP_WARNING'), this.translateService.instant('APP_SUCCESS')]
            },
            color: ['#F57660', '#F8B96A', '#90D072'],
            series: [
                {
                    name: this.translateService.instant('APP_TYPE'),
                    type: 'pie',
                    radius: ['50%', '70%'],
                    avoidLabelOverlap: false,
                    label: {
                        normal: {
                            show: false,
                            position: 'center'
                        },
                        emphasis: {
                            show: true,
                            textStyle: {
                                fontSize: '30',
                                fontWeight: 'bold'
                            }
                        }
                    },
                    labelLine: {
                        normal: {
                            show: false
                        }
                    },
                    data: []
                }
            ]
        };
        this.pieChartOptions['series'][0].data = [
            {value: item.totalSum.danger, name: this.translateService.instant('APP_DANGER')},
            {value: item.totalSum.warning, name: this.translateService.instant('APP_WARNING')},
            {value: item.totalSum.success, name: this.translateService.instant('APP_SUCCESS')},
        ];

        this.barChartOptions = {
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'shadow'
                }
            },
            legend: {
                data: [this.translateService.instant('APP_DANGER'),
                    this.translateService.instant('APP_WARNING'), this.translateService.instant('APP_SUCCESS')]
            },
            color: ['#F57660', '#F8B96A', '#90D072'],
            grid: {
                left: '3%',
                right: '4%',
                bottom: '3%',
                containLabel: true
            },
            xAxis: {
                type: 'value'
            },
            yAxis: {
                type: 'category',
                data: []
            },
            series: [
                {
                    name: this.translateService.instant('APP_DANGER'),
                    type: 'bar',
                    stack: this.translateService.instant('APP_TOTAL'),
                    label: {
                        position: 'insideRight'
                    },
                    data: []
                },
                {
                    name: this.translateService.instant('APP_WARNING'),
                    type: 'bar',
                    stack: this.translateService.instant('APP_TOTAL'),
                    label: {
                        position: 'insideRight'
                    },
                    data: []
                },
                {
                    name: this.translateService.instant('APP_SUCCESS'),
                    type: 'bar',
                    stack: this.translateService.instant('APP_TOTAL'),
                    label: {
                        position: 'insideRight'
                    },
                    data: []
                },
            ]
        };
        for (const category in item.listSum) {
            if (category) {
                this.barChartOptions.yAxis.data.push(this.translateService.instant(category));
                this.barChartOptions.series[0].data.push(item.listSum[category].danger);
                this.barChartOptions.series[1].data.push(item.listSum[category].warning);
                this.barChartOptions.series[2].data.push(item.listSum[category].success);
            }
        }
    }

    getGrade(score): string {
        if (score >= 97) {
            return 'A+';
        } else if (score >= 93) {
            return 'A';
        } else if (score >= 90) {
            return 'A-';
        } else if (score >= 87) {
            return 'B+';
        } else if (score >= 83) {
            return 'B';
        } else if (score >= 80) {
            return 'B-';
        } else if (score >= 77) {
            return 'C+';
        } else if (score >= 73) {
            return 'C';
        } else if (score >= 70) {
            return 'C-';
        } else if (score >= 67) {
            return 'D+';
        } else if (score >= 63) {
            return 'D';
        } else if (score >= 60) {
            return 'D-';
        }
        return 'F';
    }
}
