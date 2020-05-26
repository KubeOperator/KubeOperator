import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterConditionComponent } from './cluster-condition.component';

describe('ClusterConditionComponent', () => {
  let component: ClusterConditionComponent;
  let fixture: ComponentFixture<ClusterConditionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterConditionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterConditionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
