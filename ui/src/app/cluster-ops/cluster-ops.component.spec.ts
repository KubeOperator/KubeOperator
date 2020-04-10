import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterOpsComponent } from './cluster-ops.component';

describe('ClusterOpsComponent', () => {
  let component: ClusterOpsComponent;
  let fixture: ComponentFixture<ClusterOpsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterOpsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterOpsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
