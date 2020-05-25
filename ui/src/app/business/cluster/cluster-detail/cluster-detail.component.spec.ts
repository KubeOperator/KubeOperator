import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterDetailComponent } from './cluster-detail.component';

describe('ClusterDetailComponent', () => {
  let component: ClusterDetailComponent;
  let fixture: ComponentFixture<ClusterDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
