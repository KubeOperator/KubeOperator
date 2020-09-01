import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterUpgradeComponent } from './cluster-upgrade.component';

describe('ClusterUpgradeComponent', () => {
  let component: ClusterUpgradeComponent;
  let fixture: ComponentFixture<ClusterUpgradeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterUpgradeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterUpgradeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
