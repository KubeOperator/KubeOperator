import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterSettingComponent } from './multi-cluster-setting.component';

describe('MultiClusterSettingComponent', () => {
  let component: MultiClusterSettingComponent;
  let fixture: ComponentFixture<MultiClusterSettingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterSettingComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterSettingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
