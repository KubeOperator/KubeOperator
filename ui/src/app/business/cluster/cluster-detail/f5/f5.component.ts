import { Component, OnInit } from '@angular/core';
import { F5Service } from './f5.service';
import {ActivatedRoute} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {Cluster} from '../../cluster';
import {F5, F5CreateRequest} from './f5';

@Component({
  selector: 'app-f5',
  templateUrl: './f5.component.html',
  styleUrls: ['./f5.component.css']
})
export class F5Component implements OnInit {
  item: F5CreateRequest = new F5CreateRequest();
  // items: F5CreateRequest[] = [];
  currentCluster: Cluster;

  constructor(
      private f5Service: F5Service,
      private route: ActivatedRoute,
      private  http: HttpClient
  ) { }

  ngOnInit(): void {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data.cluster;
      this.f5Service.getItems(this.currentCluster.name).subscribe(d => {
        // this.item = d[0];
        this.item = d[0];
      });
    });
  }
  onSubmit() {
    window.alert('发送POST请求');
  }
  onUpdate() {
    window.alert('发送Update请求');
  }
}
