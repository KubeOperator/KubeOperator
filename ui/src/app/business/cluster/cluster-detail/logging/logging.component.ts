import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {ToolsService} from '../tools/tools.service';
import {EfComponent} from './ef/ef.component';
import {LokiComponent} from './loki/loki.component';

@Component({
    selector: 'app-logging',
    templateUrl: './logging.component.html',
    styleUrls: ['./logging.component.css']
})
export class LoggingComponent implements OnInit {

    @ViewChild(EfComponent, {static: true})
    ef: EfComponent;

    @ViewChild(LokiComponent, {static: true})
    loki: LokiComponent;

    constructor(private toolsService: ToolsService, private route: ActivatedRoute) {
    }

    loading = false;
    loggingTools = ''
    clusterName = '';
    indexPrefix = '';

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.clusterName = data.cluster.name;
            this.toolsService.list(this.clusterName).subscribe(tools => {
                for (const tool of tools) {
                    if (tool.name === 'logging' && tool.status === 'Running') {
                        this.loggingTools = 'ef'
                        this.indexPrefix = tool.vars['fluentd-elasticsearch.elasticsearch.logstashPrefix'] + '-'
                        this.loading = true;
                        break;
                    } else if (tool.name === 'loki' && tool.status === 'Running') {
                        this.loggingTools = 'loki'
                        // this.loki.currentClusterName = this.clusterName
                        this.loading = true;
                    }
                }
            })
        });
    }
    
}
