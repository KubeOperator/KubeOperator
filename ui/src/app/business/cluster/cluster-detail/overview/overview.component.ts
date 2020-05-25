import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../../cluster.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-overview',
    templateUrl: './overview.component.html',
    styleUrls: ['./overview.component.css']
})
export class OverviewComponent implements OnInit {

    constructor(private service: ClusterService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            console.log(data);
        });
    }


}
